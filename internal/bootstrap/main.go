package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord"
)

func preStart(commitHash, commitMessage string) {
	err := config.Load()
	if err != nil {
		slog.Error("Cannot load config", tint.Err(err))
		os.Exit(1)
	}

	loadGitInfo(commitHash, commitMessage)
}

func Start(commitHash, commitMessage string) {
	preStart(commitHash, commitMessage)
	setupLog(config.Config.LogType, config.Config.LogLevel)

	if err := connectToDB(); err != nil {
		slog.Error("Cannot connect to database", tint.Err(err))
		os.Exit(1)
	}

	initTracing()

	go connectToDiscord()
	go startHTTPServer()

	stop := make(chan os.Signal, 1)
	//lint:ignore SA1016 i dont know, it just works lol
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop
	err := discord.Bot.Stop()
	if err != nil {
		slog.Error("Cannot disconnect... we are disconnecting anyway...", tint.Err(err))
	}
	slog.Info("Exiting...")
}
