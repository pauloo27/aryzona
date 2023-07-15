package bot

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var SourceCommand = command.Command{
	Name: "source",
	Handler: func(ctx *command.Context) {
		t := ctx.T.(*i18n.CommandSource)
		ctx.Success(t.Description.Str(config.Config.GitRepoURL))
	},
}
