package main

import (
	"log/slog"
	"os"

	"github.com/uigraph/graphql/internal/config"
	"github.com/uigraph/graphql/internal/server"
)

func main() {
	cfg := config.Load()

	if err := server.Run(cfg); err != nil {
		slog.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
