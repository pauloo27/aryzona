package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
)

var ResumeCommand = command.Command{
	Name: "resume", Aliases: []string{"unpause"},
	Description: "Resume the queue",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		if !vc.IsPaused() {
			ctx.Error("The queue is not paused")
			return
		}
		vc.Resume()
		ctx.Success("Resumed!")
	},
}
