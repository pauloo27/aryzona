package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
)

var SkipCommand = command.Command{
	Name: "skip", Aliases: []string{"pular", "s", "sh"},
	Description: "Skip current item in the queue",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		vc.Skip()
		ctx.Success("Skipped!")
	},
}
