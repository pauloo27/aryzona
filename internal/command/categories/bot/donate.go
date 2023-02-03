package bot

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

var DonateCommand = command.Command{
	Name: "donate", Description: "Donate to the bot",
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
