package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	inventoryv1beta1 "github.com/ebisaan/proto/gateway/ebisaan/inventory/v1beta1"
	"github.com/ebisaan/proto/gateway/handler"
)

type Gateway struct {
	cfg  Config
	srv  *http.Server
	wg   sync.WaitGroup
	done chan struct{}
}

type Config struct {
	AllowedOrigins  []string
	OpenAPIFilePath string
	InventoryAddr   string
}

func New(cfg Config) *Gateway {
	return &Gateway{
		cfg:  cfg,
		done: make(chan struct{}, 1),
	}
}

func (gw *Gateway) Serve(addr string) error {
	handler, err := gw.routes()
	if err != nil {
		grpclog.Fatal(err)
	}

	gwServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	gw.srv = gwServer

	shutdownErrCh := make(chan error)
	go gw.gratefulShutdown(shutdownErrCh)

	zap.L().Info("Serving gRPC-Gateway on " + addr)
	err = gwServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrCh
	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	zap.L().Info("Stopped gRPC-Gatewayn")
	return nil
}

func (gw *Gateway) Background(fn func()) {
	gw.wg.Add(1)
	go func() {
		defer gw.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				zap.L().Error(fmt.Sprintf("Recover from: %s", err))
			}
		}()

		fn()
	}()
}

func (gw *Gateway) Done() <-chan struct{} {
	return gw.done
}

func (gw *Gateway) routes() (http.Handler, error) {
	gwMux := runtime.NewServeMux()

	conn, err := grpc.Dial(gw.config().InventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial inventory: %w", err)
	}
	ctx := context.Background()
	err = inventoryv1beta1.RegisterInventoryServiceHandler(ctx, gwMux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register inventory: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.Handle("/openapiv2", handler.OpenAPISchema(gw.config().OpenAPIFilePath))

	return handler.AllowCors(gw.config().AllowedOrigins)(mux), nil
}

func (gw *Gateway) gratefulShutdown(errCh chan<- error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	zap.L().Info(fmt.Sprintf("Receive signal %s", s))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	zap.L().Info("Shutdowning...")
	err := gw.srv.Shutdown(ctx)
	if err != nil {
		errCh <- err
	}

	zap.L().Info("Waiting for background tasks...")
	close(gw.done)
	gw.wg.Wait()
	zap.L().Info("Background tasks completed")
	close(errCh)
}

func (gw *Gateway) config() Config {
	return gw.cfg
}
