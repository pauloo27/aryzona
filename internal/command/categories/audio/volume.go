package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
)

var VolumeCommand = command.Command{
	Name:        "volume",
	Validations: []*command.Validation{validations.MustBePlaying},
	Handler: func(ctx *command.Context) command.Result {
		return ctx.ReplyRaw("https://i.imgur.com/K7v2ue7.png")
	},
}
