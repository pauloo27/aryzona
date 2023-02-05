package fun

import (
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/core/rnd"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var PickCommand = command.Command{
	Name: "pick",
	Parameters: []*command.CommandParameter{
		{
			Name:     "things",
			Required: true,
			Type:     parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandPick)

		things := strings.Split(ctx.Args[0].(string), " ")
		n, err := rnd.Rnd(len(things))
		if err != nil {
			ctx.Error(t.SomethingWentWrong.Str())
			return
		}

		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle("Pick").
				WithDescription(
					t.Description.Str(ctx.Args[0], things[n]),
				),
		)
	},
}
