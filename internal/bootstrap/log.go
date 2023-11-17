package bootstrap

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/config"
)

func setupLog(logType config.LogType, level slog.Level) {
	var handler slog.Handler

	switch logType {
	case config.LogTypeJSON:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	case config.LogTypeColored:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			TimeFormat: time.DateTime,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("Log defined", "type", logType, "level", level)
}
