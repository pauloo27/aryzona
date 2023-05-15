package animals

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/animal"
	"github.com/pauloo27/logger"
)

var CatCommand = command.Command{
	Name: "cat",
	Handler: func(ctx *command.CommandContext) {
		url, err := animal.GetRandomCat()
		if err != nil {
			ctx.Error(ctx.Lang.SomethingWentWrong.Str())
			logger.Error(err)
			return
		}
		if ctx.Reply(url) != nil {
			logger.Error(url)
		}
	},
}
