package audio

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/logger"
)

var StopCommand = command.Command{
	Name: "stop", Aliases: []string{"st", "parar", "pare"},
	Handler: func(ctx *command.CommandContext) {
		voicer := voicer.GetExistingVoicerForGuild(ctx.Message.GuildID)
		if voicer == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		err := voicer.Disconnect()
		if err != nil {
			logger.Errorf("%v", err)
			ctx.Error("Something went wrong when disconnecting...")
			return
		}
		ctx.Success("Stopped!")
	},
}
