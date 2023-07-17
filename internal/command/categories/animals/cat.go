package animals

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/animal"
	"github.com/pauloo27/logger"
)

var CatCommand = command.Command{
	Name: "cat",
	Handler: func(ctx *command.Context) command.Result {
		url, err := animal.GetRandomCat()
		if err != nil {
			logger.Error(err)
			return ctx.Error(ctx.Lang.SomethingWentWrong.Str())
		}
		return ctx.ReplyRaw(url)
	},
}
