package tools

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var PasswordCommand = command.Command{
	Name: "password", Aliases: []string{"pass"},
	Description: "Generate a random password",
	Ephemeral:   true,
	Parameters: []*command.CommandParameter{
		{Name: "length", Description: "Length of the password", Type: parameters.ParameterInt, Required: false},
	},
	Handler: func(ctx *command.CommandContext) {
		length := 10
		if len(ctx.Args) == 1 {
			length = ctx.Args[0].(int)
		}
		password, err := generatePassword(length)

		if err != nil {
			ctx.Error("Failed to generate password =(")
			return
		}

		embed := model.NewEmbed().
			WithTitle("Random password").
			WithDescription(
				fmt.Sprintf(
					"Your random password of length %d is: `%s`\nClick \"Dismiss message\" when you have it copied to delete it from Discord chat!",
					length, password,
				),
			)

		ctx.SuccessEmbed(embed)
	},
}

func generatePassword(length int) (string, error) {
	return gonanoid.New(length)
}
