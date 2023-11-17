package bootstrap

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/server"
)

func startHTTPServer() {
	if config.Config.HTTPServerPort == 0 {
		slog.Warn("HTTP server disabled")
		return
	}

	err := server.StartHTTPServer()
	if err != nil {
		slog.Error("Cannot start HTTP server", tint.Err(err))
		os.Exit(1)
	}
}
