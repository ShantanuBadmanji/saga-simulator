package cfg

import (
	"strings"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/env"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	SQLITE_DB_DIR string
	PORT          int
	LOG_DIR       string
	GO_ENV        string
}

func (e *EnvVars) isProd() bool {
	goEnv := strings.ToLower(e.GO_ENV)
	return goEnv == "prod" || goEnv == "production"
}

func (e *EnvVars) IsDev() bool {
	if e.isProd() {
		return false
	}

	return true
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

	return &EnvVars{
		SQLITE_DB_DIR: dbDir,
		PORT:          port,
		LOG_DIR:       logDir,
	}, nil
}
