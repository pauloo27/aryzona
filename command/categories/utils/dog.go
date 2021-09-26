package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/providers/animal"
	"github.com/Pauloo27/aryzona/utils"
)

var DogCommand = command.Command{
	Name: "dog", Description: "Get a cute dog",
	Aliases: []string{"woof", "doggo"},
	Handler: func(ctx *command.CommandContext) {
		url, err := animal.GetRandomDog()
		if err != nil {
			ctx.Error(utils.Fmt("An error occurred:\n %v", err))
			return
		}
		ctx.Success(url)
	},
}
