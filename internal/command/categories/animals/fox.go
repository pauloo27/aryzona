package animals

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/animal"
	"github.com/pauloo27/logger"
)

var FoxCommand = command.Command{
	Name: "fox",
	Handler: func(ctx *command.Context) {
		url, err := animal.GetRandomFox()
		if err != nil {
			ctx.Error(ctx.Lang.SomethingWentWrong.Str())
			logger.Error(err)
		}
		if ctx.Reply(url) != nil {
			logger.Error(url)
		}
	},
}
