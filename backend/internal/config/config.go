package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Server
	Port string `envconfig:"BACKEND_PORT" default:"8080"`

	// Database
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`

	// JWT
	JWTSecret        string        `envconfig:"JWT_SECRET" required:"true"`
	JWTAccessExpiry  time.Duration `envconfig:"JWT_ACCESS_EXPIRY" default:"15m"`
	JWTRefreshExpiry time.Duration `envconfig:"JWT_REFRESH_EXPIRY" default:"168h"`

	// CORS
	CORSOrigins string `envconfig:"CORS_ORIGINS" default:"http://localhost:3000"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
