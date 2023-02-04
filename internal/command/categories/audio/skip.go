package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var SkipCommand = command.Command{
	Name: "skip", Aliases: []string{"s"},
	Description: "Skip current item in the queue",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandSkip)
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		vc.Skip()
		ctx.Success(t.Skipped.Str())
	},
}
