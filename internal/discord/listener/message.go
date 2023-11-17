package listener

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/event"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

func init() {
	err := discord.Bot.Listen(event.MessageCreated, messageCreated)
	if err != nil {
		slog.Error("Cannot listen to message created event", tint.Err(err))
		os.Exit(1)
	}
}

func messageCreated(bot discord.BotAdapter, m model.Message) {
	self, err := bot.Self()
	if err != nil {
		return
	}

	if m.Author().ID() == self.ID() {
		return
	}

	if strings.HasPrefix(m.Content(), command.Prefix) {
		handleCommand(bot, self, m)
	}
}
