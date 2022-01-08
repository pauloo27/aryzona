package bootstrap

import (
	"github.com/Pauloo27/logger"
	"github.com/joho/godotenv"
)

func loadEnv() {
	logger.Info("Loading .env...")
	err := godotenv.Load()
	logger.HandleFatal(err, "Cannot load .env")
	logger.Success(".env loaded")
}
