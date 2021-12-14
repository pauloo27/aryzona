package utils

import (
	"os"

	"github.com/Pauloo27/aryzona/internal/command"
)

var SourceCommand = command.Command{
	Name: "source", Description: "Source code link",
	Aliases: []string{"s", "sauce", "src"},
	Handler: func(ctx *command.CommandContext) {
		ctx.Success("I'm a open source bot, here's my code: " + os.Getenv("DC_BOT_REMOTE_REPO"))
	},
}
