package fun

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/utils/rnd"
)

var DrawCommand = command.Command{
	Name: "draw", Description: "Draw a random thing",
	Aliases: []string{"sorteio", "sortear"},
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
			discord.NewEmbed().
				WithTitle("Draw").
				WithDescription(
					fmt.Sprintf(
						"Picking a random thing from _%s_:\n\nMy pick is **%s**",
						ctx.Args[0], things[n],
					),
				),
		)
	},
}
