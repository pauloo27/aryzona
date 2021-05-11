package utils

import (
	"github.com/Pauloo27/aryzona/command"
)

var SourceCommand = command.Command{
	Name: "source", Description: "Source code link",
	Aliases: []string{"s", "sauce", "src"},
	Handler: func(ctx *command.CommandContext) {
		ctx.Success("I'm a open source bot, here's my code: https://github.com/Pauloo27/aryzona")
	},
}
