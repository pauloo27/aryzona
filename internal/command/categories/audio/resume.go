package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var ResumeCommand = command.Command{
	Name:        "resume",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandResume)
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		if !vc.IsPaused() {
			ctx.Error(t.NotPaused.Str())
			return
		}
		vc.Resume()
		ctx.Success(t.Resumed.Str())
	},
}
