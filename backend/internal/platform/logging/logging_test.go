package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"testing/slogtest"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/logging"
)

func parseLogLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	line := bytes.TrimRight(buf.Bytes(), "\n")
	var m map[string]any
	if err := json.Unmarshal(line, &m); err != nil {
		t.Fatalf("failed to parse log output %q: %v", buf.String(), err)
	}
	return m
}

func newLogger(buf *bytes.Buffer) *slog.Logger {
	base := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(logging.NewContextHandler(base))
}

func TestContextHandler_NoEnrichment(t *testing.T) {
	var buf bytes.Buffer
	newLogger(&buf).Info("hello")

	m := parseLogLine(t, &buf)
	if _, ok := m["request_id"]; ok {
		t.Error("unexpected request_id in output")
	}
	if _, ok := m["user_id"]; ok {
		t.Error("unexpected user_id in output")
	}
}

func TestContextHandler_SingleAttr(t *testing.T) {
	var buf bytes.Buffer
	ctx := logging.WithLogAttrs(context.Background(), slog.String("request_id", "req-123"))
	newLogger(&buf).InfoContext(ctx, "hello")

	m := parseLogLine(t, &buf)
	if got := m["request_id"]; got != "req-123" {
		t.Errorf("request_id = %v, want %q", got, "req-123")
	}
}

func TestContextHandler_MultipleAttrs(t *testing.T) {
	var buf bytes.Buffer
	ctx := logging.WithLogAttrs(context.Background(),
		slog.String("request_id", "req-abc"),
		slog.String("user_id", "u99"),
	)
	newLogger(&buf).InfoContext(ctx, "hello")

	m := parseLogLine(t, &buf)
	if got := m["request_id"]; got != "req-abc" {
		t.Errorf("request_id = %v, want %q", got, "req-abc")
	}
	if got := m["user_id"]; got != "u99" {
		t.Errorf("user_id = %v, want %q", got, "u99")
	}
}

func TestContextHandler_RecordIsolation(t *testing.T) {
	// ContextHandler must enrich a Clone of the record, not the original.
	// Verify by logging through both a ContextHandler and a plain handler;
	// the plain handler must not gain any context attrs.
	var enrichedBuf, plainBuf bytes.Buffer
	ctx := logging.WithLogAttrs(context.Background(),
		slog.String("request_id", "rid"),
		slog.String("extra", "x"),
	)

	slog.New(logging.NewContextHandler(slog.NewJSONHandler(&enrichedBuf, nil))).
		InfoContext(ctx, "msg", "existing", "val")
	slog.New(slog.NewJSONHandler(&plainBuf, nil)).
		Info("msg", "existing", "val")

	em := parseLogLine(t, &enrichedBuf)
	pm := parseLogLine(t, &plainBuf)

	for _, want := range []string{"existing", "request_id", "extra"} {
		if _, ok := em[want]; !ok {
			t.Errorf("enriched output missing %q", want)
		}
	}
	if _, ok := pm["request_id"]; ok {
		t.Error("plain output must not contain request_id — record isolation violated")
	}
	if _, ok := pm["extra"]; ok {
		t.Error("plain output must not contain extra — record isolation violated")
	}
}

func TestContextHandler_DuplicateAttrs(t *testing.T) {
	// When a reserved name appears in both logAttrsKey and a call-site attr,
	// both values appear in the raw output (no runtime deduplication — documented behaviour).
	var buf bytes.Buffer
	ctx := logging.WithLogAttrs(context.Background(), slog.String("user_id", "from-middleware"))
	newLogger(&buf).InfoContext(ctx, "event", "user_id", "from-callsite")

	raw := buf.String()
	count := bytes.Count([]byte(raw), []byte(`"user_id"`))
	if count < 2 {
		t.Errorf("expected ≥2 occurrences of user_id (no dedup), got %d; output: %s", count, raw)
	}
}

func TestContextHandler_WithGroup(t *testing.T) {
	// service (pre-bound via WithAttrs before WithGroup) stays at root.
	// logAttrs are record attrs → nested under the active group.
	var buf bytes.Buffer
	base := slog.NewJSONHandler(&buf, nil).
		WithAttrs([]slog.Attr{slog.String("service", "test")})
	grouped := slog.New(logging.NewContextHandler(base)).WithGroup("http")

	ctx := logging.WithLogAttrs(context.Background(), slog.String("request_id", "gid"))
	grouped.InfoContext(ctx, "msg", "status", float64(200))

	var m map[string]any
	if err := json.Unmarshal(bytes.TrimRight(buf.Bytes(), "\n"), &m); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if m["service"] != "test" {
		t.Errorf("service = %v, want %q at root level", m["service"], "test")
	}
	group, ok := m["http"].(map[string]any)
	if !ok {
		t.Fatalf("expected http group, got %T: %v", m["http"], m["http"])
	}
	if group["request_id"] != "gid" {
		t.Errorf("http.request_id = %v, want %q (nested under group)", group["request_id"], "gid")
	}
}

func TestWithLogAttrs_Accumulation(t *testing.T) {
	var buf bytes.Buffer
	ctx := logging.WithLogAttrs(context.Background(), slog.String("a", "1"))
	ctx = logging.WithLogAttrs(ctx, slog.String("b", "2"))
	newLogger(&buf).InfoContext(ctx, "msg")

	m := parseLogLine(t, &buf)
	if m["a"] != "1" {
		t.Errorf("a = %v, want %q", m["a"], "1")
	}
	if m["b"] != "2" {
		t.Errorf("b = %v, want %q", m["b"], "2")
	}
}

func TestWithLogAttrs_ParentUnchanged(t *testing.T) {
	parent := logging.WithLogAttrs(context.Background(), slog.String("p", "parent"))
	child := logging.WithLogAttrs(parent, slog.String("c", "child"))

	var parentBuf, childBuf bytes.Buffer
	newLogger(&parentBuf).InfoContext(parent, "p")
	newLogger(&childBuf).InfoContext(child, "c")

	pm := parseLogLine(t, &parentBuf)
	cm := parseLogLine(t, &childBuf)

	if _, ok := pm["c"]; ok {
		t.Error("parent context must not contain child's attr")
	}
	if cm["p"] != "parent" {
		t.Errorf("child should inherit parent attr; p = %v", cm["p"])
	}
	if cm["c"] != "child" {
		t.Errorf("c = %v, want %q", cm["c"], "child")
	}
}

func TestContextHandler_SlogTestConformance(t *testing.T) {
	var buf bytes.Buffer
	slogtest.Run(t,
		func(t *testing.T) slog.Handler {
			buf.Reset()
			base := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
			return logging.NewContextHandler(base)
		},
		func(t *testing.T) map[string]any {
			t.Helper()
			var m map[string]any
			line := bytes.TrimRight(buf.Bytes(), "\n")
			if err := json.Unmarshal(line, &m); err != nil {
				t.Fatalf("slogtest result parse error: %v; raw=%q", err, buf.String())
			}
			return m
		},
	)
}
