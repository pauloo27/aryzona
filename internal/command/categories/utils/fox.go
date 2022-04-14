package utils

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/providers/animal"
)

var FoxCommand = command.Command{
	Name: "fox", Description: "Get a cute fox",
	Aliases: []string{"floof", "firefox", "ff"},
	Handler: func(ctx *command.CommandContext) {
		url, err := animal.GetRandomFox()
		if err != nil {
			ctx.Error(fmt.Sprintf("An error occurred:\n %v", err))
		}
		ctx.Success(url)
	},
}
