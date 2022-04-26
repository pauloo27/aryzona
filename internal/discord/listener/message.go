package listeners

import (
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/event"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/logger"
)

func init() {
	err := discord.Bot.Listen(event.MessageCreated, messageCreated)
	if err != nil {
		logger.Fatal(err)
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

	if !strings.HasPrefix(m.Content(), command.Prefix) {
		return
	}

	rawCommand := strings.TrimPrefix(strings.Split(m.Content(), " ")[0], command.Prefix)
	args := strings.Split(
		strings.TrimPrefix(strings.TrimPrefix(m.Content(), command.Prefix+rawCommand), " "), " ",
	)
	if len(args) == 1 && args[0] == "" {
		args = []string{}
	}
	event := command.Adapter{
		AuthorID: m.Author().ID(),
		GuildID:  m.Channel().Guild().ID(),
		Reply: func(ctx *command.CommandContext, msg string) error {
			_, err := discord.Bot.SendReplyMessage(m, msg)
			return err
		},
		ReplyEmbed: func(ctx *command.CommandContext, embed *discord.Embed) error {
			_, err := discord.Bot.SendReplyEmbedMessage(m, embed)
			return err
		},
	}
	command.HandleCommand(strings.ToLower(rawCommand), args, &event, bot)
}
