package listeners

import (
	"errors"
	"strings"
	"time"

	"xorm.io/xorm"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/db"
	"github.com/pauloo27/aryzona/internal/db/entity"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/event"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
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

	_, _, mentionsUser := parseCommandForPrefix(discord.AsMention(self.ID())+" ", m.Content())

	if mentionsUser {
		go handleLlama(bot, self, m)
		return
	}

	if m.Channel().Type() == model.ChannelTypeDirect {
		go handleLlama(bot, self, m)
		return
	}

	rawCommand, args, ok := parseCommand(self, m.Content())
	if !ok {
		return
	}

	var lastSentMessage model.Message
	guildID := m.Channel().Guild().ID()

	event := command.Adapter{
		MessageID: m.ID(),
		AuthorID:  m.Author().ID(),
		GuildID:   guildID,
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

	langName := getUserLanguage(m.Author().ID(), guildID)

	command.HandleCommand(
		strings.ToLower(rawCommand), args, langName, startTime, &event, bot, command.CommandTriggerMessage,
		m.Channel(),
	)
}

func getUserLanguage(userID, guildID string) i18n.LanguageName {
	var user = entity.User{ID: userID}

	hasUser, err := db.DB.Get(&user)
	if err != nil && !errors.Is(err, xorm.ErrNotExist) {
		logger.Error(err)
	}

	if hasUser {
		if user.PreferredLocale != "" {
			return user.PreferredLocale
		}

		if user.LastSlashCommandLocale != "" {
			return user.LastSlashCommandLocale
		}
	}

	if guildID == "" {
		return i18n.DefaultLanguageName
	}

	var guild = entity.Guild{ID: guildID}

	hasGuild, err := db.DB.Get(&guild)

	if err != nil && !errors.Is(err, xorm.ErrNotExist) {
		logger.Error(err)
	}

	if hasGuild {
		return guild.PreferredLocale
	}

	return i18n.DefaultLanguageName
}

func parseCommand(self model.User, content string) (rawCommand string, args []string, ok bool) {
	return parseCommandForPrefix(command.Prefix, content)
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
