package bootstrap

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/tracing"
)

func initTracing() {
	if !config.Config.Tracing.Enabled {
		slog.Info("Tracing is disabled")
		tracing.DisableTracer()
		return
	}

	err := tracing.InitTracer(
		config.Config.Tracing.Endpoint,
		config.Config.Tracing.ServiceName,
		config.Config.Env,
	)
	if err != nil {
		slog.Error("Cannot init tracer", tint.Err(err))
		os.Exit(1)
	}
	slog.Info("Tracing enabled")
}
