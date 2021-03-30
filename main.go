package main

import (
	"github.com/Pauloo27/aryzona/logger"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	logger.HandleFatal(err, "Cannot load .env")
}

func main() {
	logger.Success("Hello World")
	logger.Debug("Hello World")
	logger.Info("Hello World")
	logger.Warn("Hello World")
	logger.Error("Hello World")
}
