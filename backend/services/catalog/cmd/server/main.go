package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/httpserver"
	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/logging"
	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/requestid"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/api/rest"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/cfg"
	"github.com/ShantanuBadmanji/saga-simulator/backend/services/catalog/internal/datasource/sqlite/sqlcgen"
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
			slog.InfoContext(r.Context(), "long request canceled", "cause", r.Context().Err())
		}
	})

}

const (
	httpShutdownGracePeriod time.Duration = 10 * time.Second
)

func run(logger *slog.Logger, envVars *cfg.EnvVars) error {
	ctx := context.Background()

	mux := http.NewServeMux()
	registerBasicRoutes(mux)

	dataSourceName := filepath.Join(envVars.SQLITE_DB_DIR, "catalog.db")
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("db ping failed: %w", err)
	}
	queries := sqlcgen.New(db)

	service := rest.NewCatalogService(queries, queries, queries, logger)
	handler := rest.NewCatalogHandler(service, logger)
	handler.RegisterRoutes(mux)

	rootMux := http.NewServeMux()
	rootMux.Handle("/catalog/", http.StripPrefix("/catalog", mux))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", envVars.PORT),
		Handler: requestid.RequestIDMiddleware(rootMux),
	}

	// Cancel the context when an interrupt or termination signal is received.
	ctx, stopIntercepting := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopIntercepting()

	return httpserver.Run(ctx, server, httpShutdownGracePeriod, logger)
}

func main() {

	envVars, err := cfg.LoadEnvVars()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load environment variables: %v\n", err)
		os.Exit(1)
	}

	appCfg := cfg.NewAppConfig(envVars.IsDev())

	// moduleRoot shortens source paths in log records to be relative to the module root.
	// If go.mod cannot be found (e.g. trimmed paths in production), source paths fall
	// back to absolute — cosmetic degradation only, not a startup error.
	_, filename, _, _ := runtime.Caller(0)
	moduleRoot, _ := logging.FindModuleRoot(filename)

	logLevel := slog.LevelInfo
	if appCfg.IsDev {
		logLevel = slog.LevelDebug
	}

	logger, logCloser, err := logging.New("catalog", envVars.LOG_DIR, moduleRoot, logging.TextHandler, logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logCloser.Close()
	slog.SetDefault(logger)

	if err := run(logger, envVars); err != nil {
		logger.Error("run error", "error", err)
		os.Exit(1)
	}
}
