package fun

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/providers/joke"
	"github.com/Pauloo27/logger"
)

var JokeCommand = command.Command{
	Name: "joke", Description: "Get a random joke",
	Handler: func(ctx *command.CommandContext) {
		joke, err := joke.GetRandomJoke()
		if err != nil {
			ctx.Error("Something went wrong =(")
			logger.Error(err)
			return
		}
		jokeStr := fmt.Sprintf("%s\n\n%s", joke.Setup, joke.Punchline)
		ctx.Success(jokeStr)
	},
}
