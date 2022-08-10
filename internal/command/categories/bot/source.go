package bot

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/config"
)

var SourceCommand = command.Command{
	Name: "source", Description: "Source code link",
	Aliases: []string{"sauce", "src"},
	Handler: func(ctx *command.CommandContext) {
		ctx.Success("I'm a open source bot, here's my code: " + config.Config.GitRepoURL)
	},
}
