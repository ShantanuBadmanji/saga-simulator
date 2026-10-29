package logging

import (
	"context"
	"log/slog"
)

var _ slog.Handler = (*contextHandler)(nil)

// contextHandler enriches each log record with middleware-supplied attributes from WithLogAttrs.
//
// WithGroup behaviour: service and other pre-bound handler attrs remain at root level.
// request_id and logAttrs are appended to the record and are therefore subject to the
// inner handler's active group scope. Avoid WithGroup on the root logger; use .With for
// component scoping instead.
type contextHandler struct {
	inner slog.Handler
}

// NewContextHandler wraps inner with context-enrichment logic.
func NewContextHandler(inner slog.Handler) *contextHandler {
	return &contextHandler{inner: inner}
}

func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := logAttrsFromContext(ctx)

	if len(attrs) > 0 {
		r = r.Clone() // clone exactly once before any mutation
		r.AddAttrs(attrs...)
	}

	return h.inner.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{inner: h.inner.WithGroup(name)}
}
