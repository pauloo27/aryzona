package animals

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/animal"
	"github.com/pauloo27/logger"
)

var DogCommand = command.Command{
	Name: "dog",
	Handler: func(ctx *command.Context) command.Result {
		url, err := animal.GetRandomDog()
		if err != nil {
			logger.Error(err)
			return ctx.Error(ctx.Lang.SomethingWentWrong.Str())
		}
		return ctx.ReplyRaw(url)
	},
}
