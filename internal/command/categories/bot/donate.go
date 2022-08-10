package bot

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/discord"
)

var DonateCommand = command.Command{
	Name: "donate", Description: "Donate to the bot",
	Aliases: []string{"pix"},
	Handler: func(ctx *command.CommandContext) {
		msg := config.Config.DonateMessage
		embed := discord.NewEmbed().
			WithTitle("Donate to the bot so it continues up and running").
			WithDescription(msg)

		ctx.SuccessEmbed(embed)
	},
}
