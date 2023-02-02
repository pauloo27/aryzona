package fun

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/providers/joke"
	"github.com/Pauloo27/logger"
)

var JokeCommand = command.Command{
	Name: "joke", Description: "Get a random joke",
	Handler: func(ctx *command.CommandContext) {
		joke, err := joke.GetRandomJoke()
		if err != nil {
			ctx.Error(ctx.Lang.SomethingWentWrong.Str())
			logger.Error(err)
			return
		}
		ctx.Successf("%s\n\n%s", joke.Setup, joke.Punchline)
	},
}
