package animals

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/animal"
	"github.com/pauloo27/logger"
)

var FoxCommand = command.Command{
	Name: "fox",
	Handler: func(ctx *command.Context) command.Result {
		url, err := animal.GetRandomFox()
		if err != nil {
			logger.Error(err)
			return ctx.Error(ctx.Lang.SomethingWentWrong.Str())
		}
		return ctx.ReplyRaw(url)
	},
}
