package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/provider/animal"
	"github.com/Pauloo27/aryzona/utils"
)

var WoofCommand = command.Command{
	Name: "woof", Description: "Get a cute dog",
	Aliases: []string{"dog", "doggo"},
	Handler: func(ctx *command.CommandContext) {
		url, err := animal.GetRandomDog()
		if err != nil {
			ctx.Error(utils.Fmt("An error occured:\n %v", err))
			return
		}
		ctx.Success(url)
	},
}
