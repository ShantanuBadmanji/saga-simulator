package requestid_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/logging"
	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/requestid"
)

var hex32 = regexp.MustCompile(`^[a-f0-9]{32}$`)

func nopHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

// captureHandler stores the request seen by the downstream handler.
func captureHandler(captured **http.Request) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*captured = r
	})
}

func TestGenerateID_Length(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	requestid.RequestIDMiddleware(nopHandler()).ServeHTTP(rec, req)

	id := rec.Header().Get("X-Request-ID")
	if !hex32.MatchString(id) {
		t.Errorf("generated ID %q does not match ^[a-f0-9]{32}$", id)
	}
}

func TestRequestIDMiddleware_NoHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	requestid.RequestIDMiddleware(nopHandler()).ServeHTTP(rec, req)

	id := rec.Header().Get("X-Request-ID")
	if !hex32.MatchString(id) {
		t.Errorf("generated ID %q does not match expected format", id)
	}
}

func TestRequestIDMiddleware_ValidHeader_TrustEnabled(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "valid-id-123")

	requestid.RequestIDMiddleware(nopHandler(), requestid.TrustUpstreamID()).ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got != "valid-id-123" {
		t.Errorf("X-Request-ID = %q, want %q", got, "valid-id-123")
	}
}

func TestRequestIDMiddleware_InvalidHeader_TrustEnabled(t *testing.T) {
	for _, bad := range []string{"bad id!!", "   ", string(make([]byte, 65)), "id with\nnewline"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-ID", bad)

		requestid.RequestIDMiddleware(nopHandler(), requestid.TrustUpstreamID()).ServeHTTP(rec, req)

		id := rec.Header().Get("X-Request-ID")
		if id == bad {
			t.Errorf("invalid header %q was not rejected; got %q", bad, id)
		}
		if !hex32.MatchString(id) {
			t.Errorf("expected generated ID after rejecting %q, got %q", bad, id)
		}
	}
}

func TestRequestIDMiddleware_ValidHeader_NoTrust(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "valid-id-123")

	// No TrustUpstreamID() → always generate, ignore header.
	requestid.RequestIDMiddleware(nopHandler()).ServeHTTP(rec, req)

	got := rec.Header().Get("X-Request-ID")
	if got == "valid-id-123" {
		t.Error("header was trusted without TrustUpstreamID(); should have been ignored")
	}
	if !hex32.MatchString(got) {
		t.Errorf("expected generated ID, got %q", got)
	}
}

func TestRequestIDMiddleware_ContextPropagation(t *testing.T) {
	var captured *http.Request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	requestid.RequestIDMiddleware(captureHandler(&captured)).ServeHTTP(rec, req)

	responseID := rec.Header().Get("X-Request-ID")
	contextID := requestid.RequestIDFromContext(captured.Context())

	if responseID == "" {
		t.Fatal("X-Request-ID not set on response")
	}
	if contextID != responseID {
		t.Errorf("context ID %q != response header ID %q", contextID, responseID)
	}
}

func TestWithRequestID_RoundTrip(t *testing.T) {
	ctx := requestid.WithRequestID(context.Background(), "abc")
	if got := requestid.RequestIDFromContext(ctx); got != "abc" {
		t.Errorf("RequestIDFromContext = %q, want %q", got, "abc")
	}
}

func TestRequestIDFromContext_Missing(t *testing.T) {
	if got := requestid.RequestIDFromContext(context.Background()); got != "" {
		t.Errorf("expected empty string for missing ID, got %q", got)
	}
}

// TestRequestIDMiddleware_LogEnrichment verifies the full chain:
// RequestIDMiddleware → WithLogAttrs → ContextHandler → log output contains request_id.
func TestRequestIDMiddleware_LogEnrichment(t *testing.T) {
	var captured *http.Request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	requestid.RequestIDMiddleware(captureHandler(&captured)).ServeHTTP(rec, req)

	responseID := rec.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Fatal("X-Request-ID not set on response")
	}

	var buf bytes.Buffer
	base := slog.NewJSONHandler(&buf, nil)
	slog.New(logging.NewContextHandler(base)).InfoContext(captured.Context(), "test")

	var m map[string]any
	if err := json.Unmarshal(bytes.TrimRight(buf.Bytes(), "\n"), &m); err != nil {
		t.Fatalf("parse log: %v", err)
	}
	if m["request_id"] != responseID {
		t.Errorf("log request_id = %v, want %q", m["request_id"], responseID)
	}
}
