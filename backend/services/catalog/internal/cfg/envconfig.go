package cfg

import (
	"fmt"

	"github.com/ShantanuBadmanji/saga-simulator/backend/internal/platform/env"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	SqliteDbDir string
	Port        int
}

func LoadEnvConfig() (*EnvConfig, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}

	cfg := EnvConfig{}

	if dbDir, err := env.String(env.SqliteDbDir, nil); err != nil {
		return nil, err
	} else {
		cfg.SqliteDbDir = dbDir
	}

	if port, err := env.Int(env.Port, nil); err != nil {
		return nil, err
	} else {
		cfg.Port = port
	}

	return &cfg, nil
}
