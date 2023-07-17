package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var ResumeCommand = command.Command{
	Name:        "resume",
	Aliases:     []string{"unpause"},
	Validations: []*command.Validation{validations.MustBePlaying},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandResume)
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		if !vc.IsPaused() {
			return ctx.Error(t.NotPaused.Str())
		}
		vc.Resume()
		return ctx.Success(t.Resumed.Str())
	},
}
