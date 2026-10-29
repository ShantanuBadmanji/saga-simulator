package cfg

import (
	"os"
	"path/filepath"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/env"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	SqliteDbDir string
	Port        int
	LogDir      string
}

func LoadEnvVars() (*EnvVars, error) {
	// Optional: .env overrides nothing if vars already set in environment
	_ = godotenv.Load()

	dbDir, err := env.String(env.SqliteDbDir, nil)
	if err != nil {
		return nil, err
	}

	port, err := env.Int(env.Port, nil)
	if err != nil {
		return nil, err
	}

	defaultLog := defaultLogDir()
	logDir, err := env.String(env.LogDir, &defaultLog)
	if err != nil {
		return nil, err
	}

	return &EnvVars{
		SqliteDbDir: dbDir,
		Port:        port,
		LogDir:      logDir,
	}, nil
}

// defaultLogDir resolves logs/ relative to the running binary so the location
// is stable regardless of the process's working directory.
// Falls back to "./logs" if the binary path cannot be determined.
func defaultLogDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "logs"
	}
	// binary is at backend/out/catalog → logs should be at backend/logs
	return filepath.Clean(filepath.Join(filepath.Dir(exe), "..", "logs"))
}
