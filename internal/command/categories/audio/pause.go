package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var PauseCommand = command.Command{
	Name:        "pause",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandPause)

		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := ctx.Locals["playing"].(playable.Playable)

		if !playing.CanPause() {
			ctx.Error(t.CannotPause.Str())
			return
		}

		if vc.IsPaused() {
			ctx.Error(t.AlreadyPaused.Str(command.Prefix))
			return
		}
		vc.Pause()
		ctx.Successf(t.Paused.Str(command.Prefix))
	},
}
