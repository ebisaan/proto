package handler

import (
	"fmt"
	"net/http"
	"os"
)

func OpenAPISchema(file string) http.Handler {
	if file != "" {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("openapi schema file %s not found", file))
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, file)
		})
	} else {
		return http.NotFoundHandler()
	}
}

func AllowCors(allowOrigins []string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")

			origin := r.Header.Get("Origin")
			if origin != "" {
				for _, allow := range allowOrigins {
					if allow == "*" || origin == allow {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
							prelightHandler(w, r)
							return
						}
					}
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}

func prelightHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, PATCH, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
