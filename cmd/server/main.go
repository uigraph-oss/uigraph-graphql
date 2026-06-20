package main

import (
	"log/slog"
	"os"

	"github.com/uigraph/graphql/internal/config"
	"github.com/uigraph/graphql/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "err", err)
		os.Exit(1)
	}

	if err := server.Run(cfg); err != nil {
		slog.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
