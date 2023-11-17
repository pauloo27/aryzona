package bootstrap

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/data/db"
)

func connectToDB() error {
	slog.Info("Connecting to database...")
	conn, err := db.NewDB(&db.DBConfig{
		Host:     config.Config.DB.Host,
		Port:     config.Config.DB.Port,
		User:     config.Config.DB.User,
		Password: config.Config.DB.Password,
		Database: config.Config.DB.Database,
		SSL:      config.Config.DB.SSL,
	})
	if err != nil {
		return err
	}
	db.DB = conn

	err = conn.Ping()
	if err == nil {
		slog.Info("Connected to database")
	} else {
		slog.Error("Failed to connect to database", "err", err.Error())
		os.Exit(1)
	}

	slog.Info("Migrationg database...")
	if err = conn.Migrate(); err != nil {
		slog.Error("Failed to migrate database", tint.Err(err))
		return err
	}
	slog.Info("Database migrated")
	return nil
}
