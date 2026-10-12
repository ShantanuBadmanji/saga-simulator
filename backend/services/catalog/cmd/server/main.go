package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/httpserver"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/api/rest"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/repository"
)

func registerBasicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to Catalog Service")
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	mux.HandleFunc("GET /long-request", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(100 * time.Second):
			fmt.Fprintln(w, "ok")
		case <-r.Context().Done():
			log.Println("Long request canceled:", r.Context().Err())
		}
	})
}

const (
	httpShutdownGracePeriod time.Duration = 10 * time.Second
)

func run() error {
	mux := http.NewServeMux()
	registerBasicRoutes(mux)

	var productRepo repository.ProductRepository = nil
	service := rest.NewCatalogService(productRepo)
	handler := rest.NewCatalogHandler(service)

	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Cancel the context when an interrupt or termination signal is received.
	ctx, stopIntercepting := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopIntercepting()

	return httpserver.Run(ctx, server, httpShutdownGracePeriod)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Run Error: %v", err)
	}
}
