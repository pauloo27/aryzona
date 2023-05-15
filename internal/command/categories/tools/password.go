package tools

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var PasswordCommand = command.Command{
	Name: "password", Aliases: []string{"pass"},
	Ephemeral:   true,
	Parameters: []*command.CommandParameter{
		{Name: "length", Type: parameters.ParameterInt, Required: false},
	},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandPassword)

		length := 10
		if len(ctx.Args) == 1 {
			length = ctx.Args[0].(int)
		}
		password, err := generatePassword(length)

		if err != nil {
			ctx.Error(ctx.Lang.SomethingWentWrong.Str())
			logger.Error(err)
			return
		}

		embed := model.NewEmbed().
			WithTitle(t.Title.Str()).
			WithDescription(
				t.Description.Str(
					length, password,
				),
			)

		ctx.SuccessEmbed(embed)
	},
}

func generatePassword(length int) (string, error) {
	return gonanoid.New(length)
}
