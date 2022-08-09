package bot

import (
	"os"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
)

var DonateCommand = command.Command{
	Name: "donate", Description: "Donate to the bot",
	Aliases: []string{"pix"},
	Handler: func(ctx *command.CommandContext) {
		msg := os.Getenv("DC_BOT_DONATE_MESSAGE")
		embed := discord.NewEmbed().
			WithTitle("Donate to the bot so it continues up and running").
			WithDescription(msg)

		ctx.SuccessEmbed(embed)
	},
}
