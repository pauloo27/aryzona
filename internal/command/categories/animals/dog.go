package animals

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/providers/animal"
	"github.com/pauloo27/logger"
)

var DogCommand = command.Command{
	Name: "dog",
	Handler: func(ctx *command.Context) {
		url, err := animal.GetRandomDog()
		if err != nil {
			ctx.Error(ctx.Lang.SomethingWentWrong.Str())
			logger.Error(err)
			return
		}
		if ctx.Reply(url) != nil {
			logger.Error(err)
		}
	},
}
