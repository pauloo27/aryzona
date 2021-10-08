package audio

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord/voicer"
)

var SkipCommand = command.Command{
	Name: "skip", Aliases: []string{"pular", "sk"},
	Description: "Skip current item in the queue",
	Handler: func(ctx *command.CommandContext) {
		voicer := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if voicer == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		voicer.Skip()
		ctx.Success("Skipped!")
	},
}
