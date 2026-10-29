package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func serveHTTP(server *http.Server, logger *slog.Logger, errCh chan<- error) {
	logger.Info("starting HTTP server", "addr", server.Addr)
	errCh <- server.ListenAndServe()
	close(errCh)
}

func shutdownServer(server *http.Server, shutdownGracePeriod time.Duration, logger *slog.Logger) error {
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownGracePeriod,
	)
	defer cancel()

	logger.Info("initiating graceful shutdown")
	if err := server.Shutdown(shutdownCtx); err != nil {
		err = fmt.Errorf("failed to gracefully shut down HTTP server: %w", err)
		logger.Error("graceful shutdown failed", "error", err)

		logger.Warn("attempting to forcibly close HTTP server")
		if closeErr := server.Close(); closeErr != nil {
			closeErr = fmt.Errorf("failed to forcibly close HTTP server: %w", closeErr)
			logger.Error("force close failed", "error", closeErr)
			return errors.Join(err, closeErr)
		}

		return err
	}

	logger.Info("HTTP server shut down gracefully")
	return nil
}

func Run(ctx context.Context, server *http.Server, shutdownGracePeriod time.Duration, logger *slog.Logger) error {

	servingErrCh := make(chan error, 1)
	go serveHTTP(server, logger, servingErrCh)

	select {
	case serveErr := <-servingErrCh:
		if errors.Is(serveErr, http.ErrServerClosed) {
			logger.Info("HTTP server was already shut down or closed")
			return nil
		}
		serveErr = fmt.Errorf("HTTP server failed during startup: %w", serveErr)
		logger.Error("HTTP server failed during startup", "error", serveErr)

		logger.Info("attempting to drain remaining active requests")
		return errors.Join(serveErr, shutdownServer(server, shutdownGracePeriod, logger))
	case <-ctx.Done():
		logger.Info("signal received", "cause", context.Cause(ctx))
	}

	shutdownErr := shutdownServer(server, shutdownGracePeriod, logger)

	if serveErr := <-servingErrCh; !errors.Is(serveErr, http.ErrServerClosed) {
		return errors.Join(shutdownErr, fmt.Errorf("HTTP server failed during shutdown: %w", serveErr))
	}

	return shutdownErr
}
