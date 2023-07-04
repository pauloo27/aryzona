package bootstrap

import (
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/server"
	"github.com/pauloo27/logger"
)

func startHTTPServer() {
	if config.Config.HTTPServerPort == 0 {
		logger.Warn("HTTP server disabled")
		return
	}

	err := server.StartHTTPServer()
	if err != nil {
		logger.Fatal("Cannot start HTTP server", err)
	}
}
