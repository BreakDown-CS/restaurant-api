// Package config reads settings from environment variables and checks them at startup.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	CORSOrigins []string
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

	for _, origin := range strings.Split(os.Getenv("CORS_ORIGINS"), ",") {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		if !isOrigin(origin) {
			errs = append(errs, fmt.Errorf("CORS_ORIGINS: %q must look like http://localhost:3000 (no path, no trailing /)", origin))
			continue
		}
		cfg.CORSOrigins = append(cfg.CORSOrigins, origin)
	}

	if len(cfg.CORSOrigins) == 0 {
		errs = append(errs, errors.New("CORS_ORIGINS is required, e.g. http://localhost:3000"))
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("invalid configuration:\n%w", errors.Join(errs...))
	}

	return cfg, nil
}

func isOrigin(s string) bool {
	u, err := url.Parse(s)
	return err == nil &&
		(u.Scheme == "http" || u.Scheme == "https") &&
		u.Host != "" && u.Path == "" && u.RawQuery == ""
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
