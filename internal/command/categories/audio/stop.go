package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
)

var StopCommand = command.Command{
	Name:        "stop",
	Validations: []*command.Validation{validations.MustHaveVoicerOnGuild},
	Handler: func(ctx *command.Context) {
		t := ctx.T.(*i18n.CommandStop)
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		err := vc.Disconnect()
		if err != nil {
			ctx.Error(t.SomethingWentWrong.Str())
			logger.Error(err)
			return
		}
		ctx.Success(t.Stopped.Str())
	},
}
