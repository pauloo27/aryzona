package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/categories/sysmon"
	"github.com/Pauloo27/aryzona/command/categories/utils"
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

func registerCategory(category command.Category) {
	for _, cmd := range category.Commands {
		command.RegisterCommand(cmd)
	}
}

func main() {
	logger.Info("Connecting to Discord...")
	discord.Create(os.Getenv("DC_BOT_TOKEN"))
	discord.AddDefaultListeners()
	discord.Connect()
	logger.Success("Connected to discord")

	logger.Info("Registering commands...")
	command.Prefix = os.Getenv("DC_BOT_PREFIX")
	registerCategory(utils.Utils)
	registerCategory(sysmon.SysMon)
	logger.Success("Commands loaded")

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop
	discord.Disconnect()
	// TODO: notify before leaving
	logger.Success("Exiting...")
}
