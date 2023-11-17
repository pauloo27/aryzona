package audio

import (
	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var StopCommand = command.Command{
	Name:        "stop",
	Validations: []*command.Validation{validations.MustHaveVoicerOnGuild},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandStop)
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		err := vc.Disconnect()
		if err != nil {
			slog.Error("Cannot disconnect", tint.Err(err))
			return ctx.Error(t.SomethingWentWrong.Str())
		}
		return ctx.Success(t.Stopped.Str())
	},
}
