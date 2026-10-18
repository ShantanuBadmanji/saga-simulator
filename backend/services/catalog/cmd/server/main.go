package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/httpserver"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/api/rest"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/sqlite/sqlcgen"
	"github.com/joho/godotenv"
	_ "github.com/ncruces/go-sqlite3/driver"
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
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("Error loading .env file: %w", err)
	}

	mux := http.NewServeMux()
	registerBasicRoutes(mux)

	dbPath := filepath.Join(os.Getenv("DB_DIR"), "catalog.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("db ping failed: %w", err)
	}
	queries := sqlcgen.New(db)

	service := rest.NewCatalogService(queries, queries, queries)
	handler := rest.NewCatalogHandler(service)
	handler.RegisterRoutes(mux)

	rootMux := http.NewServeMux()
	rootMux.Handle("/catalog/", http.StripPrefix("/catalog", mux))

	server := &http.Server{
		Addr:    ":8080",
		Handler: rootMux,
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
	// Configure flags for date, time, and file location
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Add an optional text prefix to all log strings
	log.SetPrefix("[CATALOG] ")

	if err := run(); err != nil {
		log.Fatalf("Run Error: %v", err)
	}
}
