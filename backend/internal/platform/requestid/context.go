// Package requestid provides request correlation ID storage and retrieval via context.
// Reserved log field: "request_id" — only RequestIDMiddleware may set this field in logs.
package requestid

import "context"

type ctxKey struct{}

// WithRequestID stores a correlation ID in the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// RequestIDFromContext retrieves the correlation ID from the context.
// Returns an empty string if no ID has been set.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}
