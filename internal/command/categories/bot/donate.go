package bot

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var DonateCommand = command.Command{
	Name:    "donate",
	Aliases: []string{"pix"},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandDonate)

		msg := config.Config.DonateMessage
		embed := model.NewEmbed().
			WithTitle(t.Title.Str()).
			WithDescription(msg)

		ctx.SuccessEmbed(embed)
	},
}
