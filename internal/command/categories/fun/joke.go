package fun

import (
	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/joke"
)

var JokeCommand = command.Command{
	Name: "joke",
	Handler: func(ctx *command.Context) command.Result {
		joke, err := joke.GetRandomJoke()
		if err != nil {
			slog.Error("Cannot get random joke", tint.Err(err))
			return ctx.Error(ctx.Lang.SomethingWentWrong.Str())
		}
		return ctx.Successf("%s\n\n%s", joke.Setup, joke.Punchline)
	},
}
