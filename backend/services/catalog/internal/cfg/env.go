package cfg

import (
	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/env"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	SqliteDbDir string
	Port        int
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

	return &EnvVars{
		SqliteDbDir: dbDir,
		Port:        port,
	}, nil
}
