package config

import "os"

type Config struct {
	APIBaseURL string // uigraph-api base URL, e.g. http://uigraph-api:8080
	Port       string // HTTP listen port for this server
	Env        string // local | dev | prod
}

func Load() *Config {
	return &Config{
		APIBaseURL: getenv("API_BASE_URL", "http://localhost:8080"),
		Port:       getenv("PORT", "8090"),
		Env:        getenv("ENV", "local"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
