// Package platform is the glue the api runs on: config, database, HTTP server.
package platform

import (
	"errors"
	"fmt"
)

// Config is the api's runtime configuration, read once at startup.
type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

const defaultHTTPAddr = ":8080"

// LoadConfig reads the config from getenv (os.Getenv in main). DATABASE_URL is required.
func LoadConfig(getenv func(string) string) (Config, error) {
	cfg := Config{HTTPAddr: getenv("HTTP_ADDR"), DatabaseURL: getenv("DATABASE_URL")}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = defaultHTTPAddr
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("load config: %w", errors.New("DATABASE_URL is not set"))
	}
	return cfg, nil
}
