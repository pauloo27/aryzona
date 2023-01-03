package fun

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/utils/rnd"
)

var PickCommand = command.Command{
	Name: "pick", Description: "Pick a random thing",
	Aliases: []string{"sorteio", "sortear", "draw"},
	Parameters: []*command.CommandParameter{
		{
			Name:        "things",
			Description: "List of things, splitted by space",
			Required:    true,
			Type:        parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		things := strings.Split(ctx.Args[0].(string), " ")
		n, err := rnd.Rnd(len(things))
		if err != nil {
			ctx.Error("Something went wrong")
			return
		}
		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle("Pick").
				WithDescription(
					fmt.Sprintf(
						"Picking a random thing from _%s_:\n\nMy pick is **%s**",
						ctx.Args[0], things[n],
					),
				),
		)
	},
}
