package bootstrap

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/logger"

	// import arikawa implementation
	_ "github.com/pauloo27/aryzona/internal/discord/impl/arkw"

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
}
