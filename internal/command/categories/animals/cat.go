package animals

import (
	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/animal"
)

var CatCommand = command.Command{
	Name: "cat",
	Handler: func(ctx *command.Context) command.Result {
		url, err := animal.GetRandomCat()
		if err != nil {
			slog.Error("Cannot get random cat", tint.Err(err))
			return ctx.Error(ctx.Lang.SomethingWentWrong.Str())
		}
		return ctx.ReplyRaw(url)
	},
}
