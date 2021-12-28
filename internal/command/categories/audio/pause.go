package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
)

var PauseCommand = command.Command{
	Name: "pause", Aliases: []string{"resume"},
	Description: "Pause/unpause the queue",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := ctx.Locals["playing"].(playable.Playable)

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
