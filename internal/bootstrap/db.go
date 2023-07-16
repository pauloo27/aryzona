package bootstrap

import (
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/data/db"
	"github.com/pauloo27/logger"
)

func connectToDB() error {
	logger.Info("Connecting to database...")
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
		logger.Success("Connected to database")
	} else {
		logger.Fatalf("Failed to connect to database: %s", err.Error())
	}

	logger.Info("Migrationg database...")
	if err = conn.Migrate(); err != nil {
		return err
	}
	logger.Success("Database migrated")
	return nil
}
