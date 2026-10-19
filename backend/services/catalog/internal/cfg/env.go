package cfg

import (
	"fmt"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/env"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	SqliteDbDir string
	Port        int
}

func LoadEnvVars() (*EnvVars, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("load .env file: %w", err)
	}

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
