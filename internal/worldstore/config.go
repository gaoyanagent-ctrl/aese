// Package worldstore defines the isolated PostgreSQL connection contract.
package worldstore

import (
	"fmt"
	"net/url"
	"os"
)

const DatabaseURLEnv = "AESE_WORLD_DATABASE_URL"

type Config struct{ DatabaseURL string }

func FromEnv() (Config, error) {
	c := Config{DatabaseURL: os.Getenv(DatabaseURLEnv)}
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("%s is required", DatabaseURLEnv)
	}
	u, err := url.Parse(c.DatabaseURL)
	if err != nil {
		return c, fmt.Errorf("parse %s: %w", DatabaseURLEnv, err)
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return c, fmt.Errorf("%s must use postgres scheme", DatabaseURLEnv)
	}
	if u.User == nil || u.User.Username() != "aese_world_app" {
		return c, fmt.Errorf("%s must use dedicated aese_world_app role", DatabaseURLEnv)
	}
	if u.Path != "/aese_world" {
		return c, fmt.Errorf("%s must select dedicated aese_world database", DatabaseURLEnv)
	}
	return c, nil
}
