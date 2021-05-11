package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/utils"
)

var WoofCommand = command.Command{
	Name: "woof", Description: "Get a cute dog",
	Aliases: []string{"dog", "doggo"},
	Handler: func(ctx *command.CommandContext) {
		url, err := utils.GetString("https://random.dog/woof")
		if err != nil {
			ctx.Error(utils.Fmt("An error occured:\n %v", err))
		}

		ctx.Success(utils.Fmt("https://random.dog/%s", url))
	},
}
