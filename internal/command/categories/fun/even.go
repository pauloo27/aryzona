package fun

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var EvenCommand = command.Command{
	Name: "even", Aliases: []string{"odd"},
	Parameters: []*command.CommandParameter{
		{
			Name:        "number",
			Type:        parameters.ParameterInt,
			Required:    true,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandEven)

		n := ctx.Args[0].(int)
		if n&1 == 0 {
			ctx.Success(t.Even.Str())
		} else {
			ctx.Success(t.Odd.Str())
		}
	},
}
