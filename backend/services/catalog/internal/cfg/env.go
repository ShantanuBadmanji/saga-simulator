package cfg

import (
	"strings"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/env"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	SQLiteDBDir string
	Port        int
	LogDir      string
	GoEnv       string
}

func (e *EnvVars) isProd() bool {
	goEnv := strings.ToLower(e.GoEnv)
	return goEnv == "prod" || goEnv == "production"
}

func (e *EnvVars) IsDev() bool {
	return !e.isProd()

}

func LoadEnvVars() (*EnvVars, error) {
	// Optional: .env overrides nothing if vars already set in environment
	_ = godotenv.Load()

	dbDir, err := env.String(env.SQLITE_DB_DIR, nil)
	if err != nil {
		return nil, err
	}

	port, err := env.Int(env.PORT, nil)
	if err != nil {
		return nil, err
	}

	logDir, err := env.String(env.LOG_DIR, nil)
	if err != nil {
		return nil, err
	}

	defaultGoEnv := "dev"
	goEnv, err := env.String(env.GO_ENV, &defaultGoEnv)
	if err != nil {
		return nil, err
	}

	return &EnvVars{
		SQLiteDBDir: dbDir,
		Port:        port,
		LogDir:      logDir,
		GoEnv:       goEnv,
	}, nil
}
