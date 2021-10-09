package audio

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord/voicer"
)

var PauseCommand = command.Command{
	Name: "Pause", Aliases: []string{"resume"},
	Description: "Pause/unpause the queue",
	Handler: func(ctx *command.CommandContext) {
		voicer := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if voicer == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		playing := voicer.Queue.First()
		if playing == nil || !playing.CanPause() {
			ctx.Error("Cannot pause the current queue item =(")
			return
		}
		voicer.TogglePause()
		if voicer.IsPaused() {
			ctx.Success("Paused!")
			return
		}
		ctx.Success("Resumed!")
	},
}
