// Package logging provides shared structured logging for platform services.
// Reserved log fields (must not be used at call sites):
//   - "request_id" — set by requestid.RequestIDMiddleware
//   - "service"    — set by logging.New at construction time
//   - "source"     — set by slog.HandlerOptions.AddSource
//   - "user_id", "tenant_id" — reserved for future auth middleware via WithLogAttrs
package logging

import (
	"context"
	"log/slog"
)

type logAttrsKey struct{}

// WithLogAttrs appends pre-selected, safe slog.Attr values to the logging slot in context.
// Only ContextHandler reads this slot. Each call accumulates; earlier middleware attrs are preserved.
//
// Security note: slog.Attr can hold arbitrary values via slog.AnyValue. Never pass tokens,
// passwords, permissions slices, or full principal objects. Code review must enforce this.
func WithLogAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	existing, _ := ctx.Value(logAttrsKey{}).([]slog.Attr)
	merged := make([]slog.Attr, len(existing)+len(attrs))
	copy(merged, existing)
	copy(merged[len(existing):], attrs)
	return context.WithValue(ctx, logAttrsKey{}, merged)
}

func logAttrsFromContext(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(logAttrsKey{}).([]slog.Attr)
	return attrs
}
