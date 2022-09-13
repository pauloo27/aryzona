package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/logger"

	// import arikawa implementation
	_ "github.com/Pauloo27/aryzona/internal/discord/impl/arkw"

	// import listeners
	_ "github.com/Pauloo27/aryzona/internal/discord/listener"

	// import scheduler
	_ "github.com/Pauloo27/aryzona/internal/utils/scheduler"

	// import all command categories
	_ "github.com/Pauloo27/aryzona/internal/command/categories/animals"
	_ "github.com/Pauloo27/aryzona/internal/command/categories/audio"
	_ "github.com/Pauloo27/aryzona/internal/command/categories/bot"
	_ "github.com/Pauloo27/aryzona/internal/command/categories/fun"
	_ "github.com/Pauloo27/aryzona/internal/command/categories/tools"
)

func preStart(commitHash, commitMessage string) {
	_, err := config.Load()
	if err != nil {
		logger.Fatal("Cannot load config", err)
	}

	loadGitInfo(commitHash, commitMessage)
	listenToLog()
}

func Start(commitHash, commitMessage string) {
	preStart(commitHash, commitMessage)

	logger.Infof("Connecting to Discord using implementation %s...", discord.Bot.Implementation())
	err := discord.CreateBot(config.Config.Token)
	if err != nil {
		logger.Fatal(err)
	}

	err = discord.Bot.Start()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Success("Connected to discord")

	command.Prefix = config.Config.Prefix

	logger.Info("Registering slash commands handlers...")
	err = discord.Bot.RegisterSlashCommands()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Success("Slash commands created!")

	logger.Success("Up and running!")

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
