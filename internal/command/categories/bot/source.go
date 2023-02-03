package bot

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var SourceCommand = command.Command{
	Name: "source", Description: "Source code link",
	Aliases: []string{"sauce", "src"},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(i18n.CommandSource)
		ctx.Success(t.Description.Str(config.Config.GitRepoURL))
	},
}
