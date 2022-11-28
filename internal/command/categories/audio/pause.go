package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
)

var PauseCommand = command.Command{
	Name:        "pause",
	Description: "Pause the queue",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := ctx.Locals["playing"].(playable.Playable)

		if !playing.CanPause() {
			ctx.Error("Cannot pause the current queue item =(")
			return
		}

		if vc.IsPaused() {
			ctx.Error("The queue is already paused")
			return
		}
		vc.Pause()
		ctx.Successf("Paused! Use `%sresume` to resume the queue.", command.Prefix)
	},
}
