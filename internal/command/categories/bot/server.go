package bot

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/command/permissions"
	"github.com/pauloo27/aryzona/internal/db/services"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
)

var ServerCommand = command.Command{
	Name: "server",
	Parameters: []*command.Parameter{
		{
			Name: "language", Type: parameters.ParameterLowerCasedString,
			ValidValuesFunc: listValidLanguages,
			Required:        true,
		},
	},
	Permission: permissions.MustBeAdmin,
	Handler: func(ctx *command.Context) {
		t := ctx.T.(*i18n.CommandServer)

		langName := i18n.FindLanguageName(ctx.Args[0].(string))

		err := services.Guild.SetGuildOptions(ctx.GuildID, i18n.LanguageName(langName))
		if err != nil {
			ctx.Error(t.SomethingWentWrong.Str())
			logger.Error(err)
			return
		}

		ctx.Success(t.ServerOptionsChanged.Str())
	},
}
