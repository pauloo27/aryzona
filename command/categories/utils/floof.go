package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/buger/jsonparser"
)

var FloofCommand = command.Command{
	Name: "floof", Description: "Get a cute fox",
	Aliases: []string{"fox", "firefox", "ff"},
	Handler: func(ctx *command.CommandContext) {
		json, err := utils.Get("https://randomfox.ca/floof/")
		if err != nil {
			ctx.Error(utils.Fmt("An error occured:\n %v", err))
		}

		url, err := jsonparser.GetString(json, "image")
		if err != nil {
			ctx.Error(utils.Fmt("An error occured:\n %v", err))
		}

		ctx.Success(url)
	},
}
