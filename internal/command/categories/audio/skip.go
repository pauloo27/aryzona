package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
)

var SkipCommand = command.Command{
	Name: "skip", Aliases: []string{"pular", "sk"},
	Description: "Skip current item in the queue",
	Handler: func(ctx *command.CommandContext) {
		vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if vc == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		if vc.Queue.Size() == 0 {
			ctx.Error("Nothing playing...")
			return
		}
		vc.Skip()
		ctx.Success("Skipped!")
	},
}
