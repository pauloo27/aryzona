package utils

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
)

var EvenCommand = command.Command{
	Name: "even", Description: "Check if a number is even or odd",
	Parameters: []*command.CommandParameter{
		{
			Name:        "number",
			Description: "number",
			Type:        parameters.ParameterInt,
			Required:    true,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		n := ctx.Args[0].(int)
		if n&1 == 0 {
			ctx.Success("Even")
		} else {
			ctx.Success("Odd")
		}
	},
}
