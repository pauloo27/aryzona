package main

import (
	"os"
	"os/signal"

	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/logger"
	"github.com/joho/godotenv"
)

func init() {
	logger.Info("Loading .env...")
	err := godotenv.Load()
	logger.HandleFatal(err, "Cannot load .env")
	logger.Success(".env loaded")
}

func main() {
	logger.Info("Connecting to Discord...")
	discord.Create(os.Getenv("DC_BOT_TOKEN"))
	discord.AddDefaultListeners()
	discord.Connect()
	logger.Success("Connected to discord")

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	logger.Success("Exiting...")
}
