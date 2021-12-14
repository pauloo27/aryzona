package utils

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/providers/xkcd"
	"github.com/Pauloo27/aryzona/internal/utils"
)

var XkcdCommand = command.Command{
	Name: "xkcd", Description: "Get the latest XKCD",
	Aliases: []string{},
	Handler: func(ctx *command.CommandContext) {
		url, err := xkcd.GetLatest()
		if err != nil {
			ctx.Error(utils.Fmt("An error occurred:\n %v", err))
			return
		}
		ctx.Success(url)
	},
}
