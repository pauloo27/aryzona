package utils

import (
	"github.com/Pauloo27/aryzona/command"
)

var EvenCommand = command.Command{
	Name: "even", Description: "Check if a number is even or odd",
	Arguments: []*command.CommandArgument{
		{
			Name:        "number",
			Description: "number",
			Type:        command.ArgumentInt,
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
