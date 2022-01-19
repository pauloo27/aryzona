package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
)

var ShuffleCommand = command.Command{
	Name:        "shuffle",
	Aliases:     []string{"s"},
	Description: "Shuffle current queue",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		vc.Queue.Shuffle()
		ctx.Success("Shuffled! :wink:")
	},
}
