package fun

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/core/f"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/providers/xkcd"
	"github.com/Pauloo27/logger"
)

var XkcdCommand = command.Command{
	Name: "xkcd",
	SubCommands: []*command.Command{
		&XkcdLatestSubCommand,
		&XkcdRandomSubCommand,
		&XkcdNumberSubCommand,
	},
}

var XkcdLatestSubCommand = command.Command{
	Name: "latest", 
	Handler: func(ctx *command.CommandContext) {
		comic, err := xkcd.GetLatest()
		sendComic(ctx, comic, err)
	},
}

var XkcdRandomSubCommand = command.Command{
	Name: "random", 
	Handler: func(ctx *command.CommandContext) {
		comic, err := xkcd.GetRandom()
		sendComic(ctx, comic, err)
	},
}

var XkcdNumberSubCommand = command.Command{
	Name: "number", 
	Aliases: []string{"num"},
	Parameters: []*command.CommandParameter{
		{
			Name: "number", 
			Type: parameters.ParameterInt, Required: true,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		comic, err := xkcd.GetByNum(ctx.Args[0].(int))
		sendComic(ctx, comic, err)
	},
}

func sendComic(ctx *command.CommandContext, comic *xkcd.Comic, err error) {
	if err != nil {
		ctx.Error(ctx.Lang.SomethingWentWrong.Str())
		logger.Error(err)
		return
	}

	ctx.SuccessEmbed(
		model.NewEmbed().
			WithTitle(fmt.Sprintf(
				"#%d - %s (%s/%s/%s)", comic.Num, comic.SafeTitle,
				// FIXME: i18n the date format
				comic.Year, f.PadLeft(comic.Month, "0", 2), f.PadLeft(comic.Day, "0", 2)),
			).
			WithURL(fmt.Sprintf("https://www.explainxkcd.com/wiki/index.php/%d", comic.Num)).
			WithImage(comic.Img).
			WithFooter(comic.Alt),
	)
}
