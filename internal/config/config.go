package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	APIBaseURL     string   // uigraph-api base URL, e.g. http://uigraph-api:8080
	Port           string   // HTTP listen port for this server
	Env            string   // local | dev | prod
	AllowedOrigins []string // CORS allow-list; empty disables CORS handling entirely
}

func Load() (*Config, error) {
	cfg := &Config{
		APIBaseURL: getenv("API_BASE_URL", "http://localhost:8080"),
		Port:       getenv("PORT", "8090"),
		Env:        getenv("ENV", "local"),
	}
	if v := getenv("ALLOWED_ORIGINS", ""); v != "" {
		cfg.AllowedOrigins = strings.Split(v, ",")
	}

	u, err := url.Parse(cfg.APIBaseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("API_BASE_URL %q is not a valid absolute URL", cfg.APIBaseURL)
	}

	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
