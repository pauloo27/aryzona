package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var ShuffleCommand = command.Command{
	Name:        "shuffle",
	Aliases:     []string{"sh"},
	Validations: []*command.Validation{validations.MustBePlaying},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandShuffle)
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		vc.Queue.Shuffle()
		return ctx.Success(t.Shuffled.Str(":wink:"))
	},
}
