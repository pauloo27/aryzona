package listeners

import (
	"strings"
	"time"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/event"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/logger"
)

func init() {
	err := discord.Bot.Listen(event.MessageCreated, messageCreated)
	if err != nil {
		logger.Fatal(err)
	}
}

func messageCreated(bot discord.BotAdapter, m model.Message) {
	startTime := time.Now()
	self, err := bot.Self()
	if err != nil {
		return
	}

	if m.Author().ID() == self.ID() {
		return
	}

	rawCommand, args, ok := parseCommand(self, m.Content())
	if !ok {
		return
	}

	var lastSentMessage model.Message

	event := command.Adapter{
		AuthorID: m.Author().ID(),
		GuildID:  m.Channel().Guild().ID(),
		Reply: func(ctx *command.CommandContext, msg string) error {
			var err error
			lastSentMessage, err = discord.Bot.SendReplyMessage(m, msg)
			return err
		},
		ReplyEmbed: func(ctx *command.CommandContext, embed *model.Embed) error {
			var err error
			lastSentMessage, err = discord.Bot.SendReplyEmbedMessage(m, embed)
			return err
		},
		ReplyComplex: func(ctx *command.CommandContext, message *model.ComplexMessage) error {
			var err error
			if message.ReplyTo == nil {
				message.ReplyTo = m
			}
			lastSentMessage, err = discord.Bot.SendComplexMessage(ctx.Channel.ID(), message)
			return err
		},
		Edit: func(ctx *command.CommandContext, msg string) error {
			_, err := discord.Bot.EditMessageContent(lastSentMessage, msg)
			return err
		},
		EditEmbed: func(ctx *command.CommandContext, embed *model.Embed) error {
			_, err := discord.Bot.EditMessageEmbed(lastSentMessage, embed)
			return err
		},
		EditComplex: func(ctx *command.CommandContext, message *model.ComplexMessage) error {
			_, err := discord.Bot.EditComplexMessage(lastSentMessage, message)
			return err
		},
	}
	command.HandleCommand(
		strings.ToLower(rawCommand), args, i18n.DefaultLanguageName, startTime, &event, bot, command.CommandTriggerMessage,
		m.Channel(),
	)
}

func parseCommand(self model.User, content string) (rawCommand string, args []string, ok bool) {
	rawCommand, args, ok = parseCommandForPrefix(command.Prefix, content)
	if ok {
		return
	}
	// check for "@bot command"
	return parseCommandForPrefix(discord.AsMention(self.ID())+" ", content)
}

func parseCommandForPrefix(prefix string, content string) (rawCommand string, args []string, ok bool) {
	if !strings.HasPrefix(content, prefix) {
		return
	}

	ok = true
	parts := strings.Split(strings.TrimPrefix(content, prefix), " ")
	rawCommand = parts[0]
	args = parts[1:]
	return
}
