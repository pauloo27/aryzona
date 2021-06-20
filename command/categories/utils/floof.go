package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/providers/animal"
	"github.com/Pauloo27/aryzona/utils"
)

var FloofCommand = command.Command{
	Name: "floof", Description: "Get a cute fox",
	Aliases: []string{"fox", "firefox", "ff"},
	Handler: func(ctx *command.CommandContext) {
		url, err := animal.GetRandomFox()
		if err != nil {
			ctx.Error(utils.Fmt("An error occured:\n %v", err))
		}
		ctx.Success(url)
	},
}
