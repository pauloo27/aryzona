package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/providers/animal"
	"github.com/Pauloo27/aryzona/utils"
)

var CatCommand = command.Command{
	Name: "cat", Description: "Get a cute cat",
	Aliases: []string{"meow", "miau"},
	Handler: func(ctx *command.CommandContext) {
		url, err := animal.GetRandomCat()
		if err != nil {
			ctx.Error(utils.Fmt("An error occurred:\n %v", err))
			return
		}
		ctx.Success(url)
	},
}
