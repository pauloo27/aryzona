package fun

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/joke"
	"github.com/pauloo27/logger"
)

var JokeCommand = command.Command{
	Name: "joke",
	Handler: func(ctx *command.Context) command.Result {
		joke, err := joke.GetRandomJoke()
		if err != nil {
			logger.Error(err)
			return ctx.Error(ctx.Lang.SomethingWentWrong.Str())
		}
		return ctx.Successf("%s\n\n%s", joke.Setup, joke.Punchline)
	},
}
