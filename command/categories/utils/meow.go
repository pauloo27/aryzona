package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/buger/jsonparser"
)

var MeowCommand = command.Command{
	Name: "meow", Description: "Get a cute cat",
	Aliases: []string{"cat", "miau"},
	Handler: func(ctx *command.CommandContext) {
		json, err := utils.Get("https://aws.random.cat/meow")
		if err != nil {
			ctx.Error(utils.Fmt("An error occured:\n %v", err))
		}

		url, err := jsonparser.GetString(json, "file")
		if err != nil {
			ctx.Error(utils.Fmt("An error occured:\n %v", err))
		}

		ctx.Success(url)
	},
}
