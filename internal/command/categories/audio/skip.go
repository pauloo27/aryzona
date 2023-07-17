package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var SkipCommand = command.Command{
	Name: "skip", Aliases: []string{"s"},
	Validations: []*command.Validation{validations.MustBePlaying},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandSkip)
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		vc.Skip()
		return ctx.Success(t.Skipped.Str())
	},
}
