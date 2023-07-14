package listeners

import (
	"strings"

	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/providers/llama"
	"github.com/pauloo27/logger"
)

func handleLlama(bot discord.BotAdapter, self model.User, m model.Message) {
	err := bot.StartTyping(m.Channel())
	if err != nil {
		logger.Error(err)
		return
	}

	msg := strings.TrimPrefix(m.Content(), discord.AsMention(self.ID())+" ")
	response, err := llama.AskLlama(nil, msg)
	if err != nil {
		logger.Error(err)
		return
	}
	response = strings.TrimPrefix(response, "Aryzona: ")
	_, err = bot.SendReplyMessage(m, response)
	if err != nil {
		logger.Error(err)
	}
}
