package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/categories/audio"
	"github.com/Pauloo27/aryzona/command/categories/sysmon"
	"github.com/Pauloo27/aryzona/command/categories/utils"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/git"
	"github.com/Pauloo27/aryzona/logger"
	"github.com/joho/godotenv"
)

var commitHash, commitMessage string

func init() {
	logger.Info("Loading .env...")
	err := godotenv.Load()
	logger.HandleFatal(err, "Cannot load .env")
	logger.Success(".env loaded")

	git.CommitHash = commitHash
	git.CommitMessage = commitMessage
	git.RemoteRepo = os.Getenv("DC_BOT_REMOTE_REPO")
}

func registerCategory(category command.Category) {
	if category.OnLoad != nil {
		category.OnLoad()
	}
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
	registerCategory(audio.Audio)
	logger.Success("Commands loaded")

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop
	discord.Disconnect()
	logger.Success("Exiting...")
}
