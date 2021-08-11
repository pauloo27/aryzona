package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/providers/animal"
	"github.com/Pauloo27/aryzona/utils"
)

var MeowCommand = command.Command{
	Name: "meow", Description: "Get a cute cat",
	Aliases: []string{"cat", "miau"},
	Handler: func(ctx *command.CommandContext) {
		url, err := animal.GetRandomCat()
		if err != nil {
			ctx.Error(utils.Fmt("An error occurred:\n %v", err))
			return
		}
		ctx.Success(url)
	},
}
