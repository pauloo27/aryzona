package bot

import (
	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/command/permissions"
	"github.com/pauloo27/aryzona/internal/data/services"
	"github.com/pauloo27/aryzona/internal/i18n"
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
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandServer)

		langName := i18n.FindLanguageName(ctx.Args[0].(string))

		err := services.Guild.SetGuildOptions(ctx.GuildID, i18n.LanguageName(langName))
		if err != nil {
			slog.Error("Cannot set guild options", tint.Err(err))
			return ctx.Error(t.SomethingWentWrong.Str())
		}

		return ctx.Success(t.ServerOptionsChanged.Str())
	},
}
