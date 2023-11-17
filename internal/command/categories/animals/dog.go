package animals

import (
	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/animal"
)

var DogCommand = command.Command{
	Name: "dog",
	Handler: func(ctx *command.Context) command.Result {
		url, err := animal.GetRandomDog()
		if err != nil {
			slog.Error("Cannot get random dog", tint.Err(err))
			return ctx.Error(ctx.Lang.SomethingWentWrong.Str())
		}
		return ctx.ReplyRaw(url)
	},
}
