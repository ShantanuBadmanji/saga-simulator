package requestlogger

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriterWrapper initializes the wrapper with a default 200 OK status.
func NewResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default if WriteHeader is never called explicitly
	}
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func RequestLoggerMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapper := NewResponseWriterWrapper(w)
		next.ServeHTTP(wrapper, r)
		duration := time.Since(startTime)
		logRequest(wrapper, r, duration, logger)
	})
}

func logRequest(w *responseWriterWrapper, r *http.Request, duration time.Duration, logger *slog.Logger) {
	logger.Info("http request processed",
		slog.Group("http",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.Int("status", w.statusCode),
			slog.Duration("duration", duration),
		),
	)
}
