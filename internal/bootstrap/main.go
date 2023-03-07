package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/server"
	"github.com/Pauloo27/logger"
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

	if err := connectToDB(); err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	go connectToDiscord()
	go server.StartHTTPServer()

	stop := make(chan os.Signal, 1)
	//lint:ignore SA1016 i dont know, it just works lol
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop
	err := discord.Bot.Stop()
	if err != nil {
		logger.Error("Cannot disconnect... we are disconnecting anyway...", err)
	}
	logger.Success("Exiting...")
}
