package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/utils"
)

var PingCommand = command.Command{
	Name: "ping", Description: "Get the bot latency",
	Aliases: []string{"pong"},
	Handler: func(ctx *command.CommandContext) {
		ctx.SuccessEmbed(
			discord.NewEmbed().
				WithDescription(
					utils.Fmt("Latency with Discord's server is **%d ms**",
						ctx.Bot.Latency().Milliseconds()),
				),
		)
	},
}
