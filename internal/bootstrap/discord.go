package bootstrap

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord"

	// disgo impl
	_ "github.com/pauloo27/aryzona/internal/discord/impl/disgo"

	// import listeners
	_ "github.com/pauloo27/aryzona/internal/discord/listener"

	// import scheduler
	_ "github.com/pauloo27/aryzona/internal/core/scheduler"

	// import all command categories
	_ "github.com/pauloo27/aryzona/internal/command/categories/animals"
	_ "github.com/pauloo27/aryzona/internal/command/categories/audio"
	_ "github.com/pauloo27/aryzona/internal/command/categories/bot"
	_ "github.com/pauloo27/aryzona/internal/command/categories/fun"
	_ "github.com/pauloo27/aryzona/internal/command/categories/tools"
)

func connectToDiscord() {
	slog.Info("Connecting to Discord", "implementation", discord.Bot.Implementation())
	err := discord.CreateBot(config.Config.Token)
	if err != nil {
		slog.Error("Cannot create bot", tint.Err(err))
		os.Exit(1)
	}

	err = discord.Bot.Start()
	if err != nil {
		slog.Error("Cannot start bot", tint.Err(err))
		os.Exit(1)
	}
	slog.Info("Connected to discord")

	command.Prefix = config.Config.Prefix

	slog.Info("Registering slash commands handlers...")
	err = discord.Bot.RegisterSlashCommands()
	if err != nil {
		slog.Error("Cannot register slash commands", tint.Err(err))
		os.Exit(1)
	}
	slog.Info("Slash commands created!")

	slog.Info("Up and running!")
}
