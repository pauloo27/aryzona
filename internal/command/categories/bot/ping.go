package bot

import (
	"fmt"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var PingCommand = command.Command{
	Name: "ping",
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandPing)

		latency := formatAPILatency(ctx.Bot)
		if latency == "0" {
			latency = t.StillCalculating.Str(":hourglass_flowing_sand:")
		}

		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle(t.Title.Str(":ping_pong:")).
				WithFooter(t.Footer.Str()).
				WithField(
					t.APILatency.Str(),
					latency,
				),
		)
	},
}

func formatAPILatency(bot discord.BotAdapter) string {
	latency := bot.Latency()
	if latency == 0 {
		return "0"
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
