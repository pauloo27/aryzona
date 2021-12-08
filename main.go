package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/git"
	"github.com/Pauloo27/logger"
	"github.com/joho/godotenv"

	// import discordgo implementation
	_ "github.com/Pauloo27/aryzona/discord/impl/dcgo"

	// import listeners
	_ "github.com/Pauloo27/aryzona/discord/listener"

	// import scheduler
	_ "github.com/Pauloo27/aryzona/utils/scheduler"

	// import all command categories
	_ "github.com/Pauloo27/aryzona/command/categories/sysmon"
	_ "github.com/Pauloo27/aryzona/command/categories/utils"
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

func main() {
	logger.Info("Connecting to Discord...")
	err := discord.CreateBot(os.Getenv("DC_BOT_TOKEN"))
	if err != nil {
		logger.Fatal(err)
	}

	err = discord.Bot.Start()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Success("Connected to discord")

	command.Prefix = os.Getenv("DC_BOT_PREFIX")

	stop := make(chan os.Signal, 1)
	//lint:ignore SA1016 i dont know, it just works lol
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop
	err = discord.Bot.Stop()
	if err != nil {
		logger.Error("Cannot disconnect... we are disconnecting anyway...", err)
	}
	logger.Success("Exiting...")
}
