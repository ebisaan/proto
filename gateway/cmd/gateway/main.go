package main

import (
	"flag"
	"log"
	"strings"

	"go.uber.org/zap"

	"github.com/ebisaan/proto/gateway/internal/gateway"
	"github.com/ebisaan/proto/gateway/internal/logger"
)

var (
	addrFlag           = flag.String("addr", ":8080", "Address to run on")
	openAPIFileFlag    = flag.String("openapi-file", "doc/ebisaan.swagger.json", "OpenAPI file path")
	allowedOriginsFlag = flag.String("allowed-origins", "", "Comma-separated string of CORS origin")
	inventoryAddrFlag  = flag.String("inventory-addr", ":8081", "Inventory service's address")
	logFormatAddrFlag  = flag.String("log-format", "json", "Log's format (json or console)")
)

func main() {
	flag.Parse()
	setLogger(*logFormatAddrFlag)

	allowedOrigins := strings.Split(*allowedOriginsFlag, ",")
	gwConfig := gateway.Config{
		AllowedOrigins:  allowedOrigins,
		OpenAPIFilePath: *openAPIFileFlag,
		InventoryAddr:   *inventoryAddrFlag,
	}

	gw := gateway.New(gwConfig)
	err := gw.Serve(*addrFlag)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}

func setLogger(format string) {
	lgr, err := logger.New(zap.InfoLevel, format)
	if err != nil {
		log.Fatal(err)
	}

	zap.ReplaceGlobals(lgr)
}
