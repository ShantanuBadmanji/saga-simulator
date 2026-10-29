package requestid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/logging"
)

type middlewareConfig struct {
	trustUpstreamID bool
}

// MiddlewareOption configures RequestIDMiddleware.
type MiddlewareOption func(*middlewareConfig)

// TrustUpstreamID instructs the middleware to accept a validated X-Request-ID header
// from the upstream request. Use this only when the service is behind a trusted gateway
// that strips or overwrites client-supplied correlation headers.
// Default: always generate a new ID regardless of incoming headers.
func TrustUpstreamID() MiddlewareOption {
	return func(c *middlewareConfig) { c.trustUpstreamID = true }
}

var validRequestID = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// RequestIDMiddleware injects a correlation ID into the context and sets X-Request-ID
// on the response. By default it always generates a new ID. Pass TrustUpstreamID()
// when deployed behind a trusted gateway.
func RequestIDMiddleware(next http.Handler, opts ...MiddlewareOption) http.Handler {
	cfg := &middlewareConfig{}
	for _, o := range opts {
		o(cfg)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ""
		if cfg.trustUpstreamID {
			candidate := r.Header.Get("X-Request-ID")
			if validRequestID.MatchString(candidate) {
				id = candidate
			}
		}
		if id == "" {
			var err error
			id, err = generateID()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
		ctx := WithRequestID(r.Context(), id)
		ctx = logging.WithLogAttrs(ctx, slog.String("request_id", id))

		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("requestid: failed to generate ID: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}
