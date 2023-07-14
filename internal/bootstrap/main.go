package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/providers/llama"
	"github.com/pauloo27/logger"
)

func preStart(commitHash, commitMessage string) {
	err := config.Load()
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
	go startHTTPServer()
	go func() {
		err := llama.Init()
		if err == nil {
			logger.Success("Llama prompt tokenized")
		} else {
			logger.Fatal(err, "Cannot tokenize prompt")
		}
	}()

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
