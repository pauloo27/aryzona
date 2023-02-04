package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var StopCommand = command.Command{
	Name: "stop", Aliases: []string{"st", "parar", "pare"},
	Description: "Stop what is playing",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandStop)
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		err := vc.Disconnect()
		if err != nil {
			ctx.Error(t.SomethingWentWrong.Str())
			return
		}
		ctx.Success(t.Stopped.Str())
	},
}
