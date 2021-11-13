package audio

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord/voicer"
)

var PauseCommand = command.Command{
	Name: "Pause", Aliases: []string{"resume"},
	Description: "Pause/unpause the queue",
	Handler: func(ctx *command.CommandContext) {
		vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if vc == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		playing := vc.Queue.First()
		if playing == nil {
			ctx.Error("Nothing playing...")
			return
		}
		if !playing.CanPause() {
			ctx.Error("Cannot pause the current queue item =(")
			return
		}
		vc.TogglePause()
		if vc.IsPaused() {
			ctx.Success("Paused!")
			return
		}
		ctx.Success("Resumed!")
	},
}
