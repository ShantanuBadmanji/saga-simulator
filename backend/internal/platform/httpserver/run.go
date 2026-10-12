package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func serveHTTP(server *http.Server, errCh chan<- error) {
	log.Println("Starting HTTP server on", server.Addr)
	errCh <- server.ListenAndServe() // not a blocking send, since the channel has a buffer of 1
	close(errCh)
}

func shutdownServer(server *http.Server, shutdownGracePeriod time.Duration) error {
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownGracePeriod,
	)
	defer cancel()

	log.Printf("Initiating graceful shutdown")
	if err := server.Shutdown(shutdownCtx); err != nil {
		err = fmt.Errorf("failed to gracefully shut down HTTP server: %w", err)
		log.Println(err)

		log.Println("Attempting to forcibly close HTTP server")
		if closeErr := server.Close(); closeErr != nil {
			closeErr = fmt.Errorf("failed to forcibly close HTTP server: %w", closeErr)
			log.Println(closeErr)
			return errors.Join(err, closeErr)
		}

		return err
	}

	log.Println("HTTP server shut down gracefully")
	return nil
}

func Run(ctx context.Context, server *http.Server, shutdownGracePeriod time.Duration) error {

	servingErrCh := make(chan error, 1)
	go serveHTTP(server, servingErrCh)

	select {
	case serveErr := <-servingErrCh:
		if errors.Is(serveErr, http.ErrServerClosed) {
			log.Println("HTTP server was already shut down or closed.")
			return nil
		}
		serveErr = fmt.Errorf("HTTP server failed during startup: %w", serveErr)
		log.Println(serveErr)

		log.Println("Attempting to drain any remaining active requests.")
		return errors.Join(serveErr, shutdownServer(server, shutdownGracePeriod))
	case <-ctx.Done():
		log.Printf("Signal interception loop ended. Cause: %v", context.Cause(ctx))
	}

	shutdownErr := shutdownServer(server, shutdownGracePeriod)

	if serveErr := <-servingErrCh; !errors.Is(serveErr, http.ErrServerClosed) {
		return errors.Join(shutdownErr, fmt.Errorf("HTTP server failed during shutdown: %w", serveErr))
	}

	return shutdownErr
}
