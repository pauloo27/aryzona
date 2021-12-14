package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/logger"
)

var StopCommand = command.Command{
	Name: "stop", Aliases: []string{"st", "parar", "pare"},
	Description: "Stop what is playing",
	Handler: func(ctx *command.CommandContext) {
		voicer := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if voicer == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		err := voicer.Disconnect()
		if err != nil {
			logger.Error(err)
			ctx.Error("Something went wrong when disconnecting...")
			return
		}
		ctx.Success("Stopped!")
	},
}
