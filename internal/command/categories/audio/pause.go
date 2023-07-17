package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var PauseCommand = command.Command{
	Name:        "pause",
	Validations: []*command.Validation{validations.MustBePlaying},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandPause)

		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := ctx.Locals["playing"].(playable.Playable)

		if !playing.CanPause() {
			return ctx.Error(t.CannotPause.Str())
		}

		if vc.IsPaused() {
			return ctx.Error(t.AlreadyPaused.Str(command.Prefix))
		}
		vc.Pause()
		return ctx.Successf(t.Paused.Str(command.Prefix))
	},
}
