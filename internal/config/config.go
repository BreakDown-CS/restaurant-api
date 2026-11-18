package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
}

func (c *Config) IsProduction() bool { return c.Env == "production" }

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Env:         getenv("APP_ENV", "development"),
		Port:        getenv("APP_PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	var errs []error
	switch cfg.Env {
	case "development", "test", "production":
	default:
		errs = append(errs, fmt.Errorf("APP_ENV must be development, test or production (got %q)", cfg.Env))
	}

	if port, err := strconv.Atoi(cfg.Port); err != nil || port < 1 || port > 65535 {
		errs = append(errs, fmt.Errorf("APP_PORT must be a number from 1 to 65535 (got %q)", cfg.Port))
	}

	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("invalid configuration:\n%w", errors.Join(errs...))
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
