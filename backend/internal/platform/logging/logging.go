package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// FindModuleRoot walks up from the given source file path until it finds a go.mod,
// returning that directory as the module root. Safe to call with runtime.Caller(0).
func FindModuleRoot(sourceFile string) (string, error) {
	dir := filepath.Dir(sourceFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", sourceFile)
		}
		dir = parent
	}
}

// HandlerType selects the log output format.
type HandlerType uint8

const (
	TextHandler HandlerType = iota + 1
	JSONHandler
)

// New creates a structured logger that writes to both stdout and {logDir}/{serviceName}.log.
// It formats timestamps as RFC3339, adds source locations relative to moduleRoot (pass empty string
// to use absolute paths), and enriches records with context values via ContextHandler.
// The caller must close the returned io.Closer on shutdown.
func New(serviceName, logDir, moduleRoot string, handlerType HandlerType, level slog.Level) (*slog.Logger, io.Closer, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, nil, err
	}
	path := filepath.Join(logDir, serviceName+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			case slog.SourceKey:
				if source, ok := a.Value.Any().(*slog.Source); ok {
					a.Value = slog.StringValue(formatSourcePath(source, moduleRoot))
				}
			}
			return a
		},
		Level: level,
	}

	var stdoutHandler slog.Handler
	switch handlerType {
	case TextHandler:
		stdoutHandler = slog.NewTextHandler(os.Stdout, opts)
	case JSONHandler:
		stdoutHandler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		f.Close()
		return nil, nil, fmt.Errorf("unsupported handler type: %v", handlerType)
	}

	mutiHandler := slog.NewMultiHandler(stdoutHandler, slog.NewJSONHandler(f, opts))
	handler := mutiHandler.WithAttrs([]slog.Attr{slog.String("service", serviceName)})
	return slog.New(NewContextHandler(handler)), f, nil
}

func formatSourcePath(source *slog.Source, moduleRoot string) string {
	if source == nil {
		return ""
	}
	if moduleRoot == "" {
		return fmt.Sprintf("%s:%d", source.File, source.Line)
	}
	rel, err := filepath.Rel(moduleRoot, source.File)
	if err != nil {
		return fmt.Sprintf("%s:%d", source.File, source.Line)
	}
	return fmt.Sprintf("%s:%d", rel, source.Line)
}
