package bot

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"
)

var PingCommand = command.Command{
	Name: "ping", Description: "Get the bot latency",
	Aliases: []string{"pong"},
	Handler: func(ctx *command.CommandContext) {
		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle(":ping_pong: Pong!").
				WithFooter("(that's the Bot latency, not yours)").
				WithField(
					"API Latency",
					formatAPILatency(ctx.Bot),
				),
		)
	},
}

func formatAPILatency(bot discord.BotAdapter) string {
	latency := bot.Latency()
	if latency == 0 {
		return "🤔 I'm still calculating..."
	}
	ms := latency.Milliseconds()
	var icon string
	if ms < 50 {
		icon = "🟢"
	} else if ms < 100 {
		icon = "🟡"
	} else {
		icon = "🔴"
	}
	return fmt.Sprintf(
		"%s %d ms",
		icon,
		ms,
	)
}
