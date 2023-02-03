package animals

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/providers/animal"
	"github.com/Pauloo27/logger"
)

var CatCommand = command.Command{
	Name: "cat", Description: "Get a cute cat",
	Aliases: []string{"meow", "miau"},
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
