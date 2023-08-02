package listener

import (
	"strings"
	"time"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

func handleCommand(bot discord.BotAdapter, self model.User, m model.Message) {
	eventTime := time.Now()
	var lastSentMessage model.Message
	guildID := m.Channel().Guild().ID()

	commandName, args, ok := parseCommand(self, m.Content())
	if !ok {
		return
	}

	trigger := command.TriggerEvent{
		Type:      command.CommandTriggerMessage,
		EventTime: eventTime,
		MessageID: m.ID(),
		AuthorID:  m.Author().ID(),
		GuildID:   guildID,
		Channel:   m.Channel(),
		Reply: func(ctx *command.Context, message *model.ComplexMessage) error {
			var err error
			if message.ReplyTo == nil {
				message.ReplyTo = m
			}
			lastSentMessage, err = discord.Bot.SendComplexMessage(ctx.Channel.ID(), message)
			return err
		},
		Edit: func(ctx *command.Context, message *model.ComplexMessage) error {
			_, err := discord.Bot.EditComplexMessage(lastSentMessage, message)
			return err
		},
	}

	cmd := command.GetCommand(commandName)
	if cmd == nil {
		return
	}

	command.HandleCommand(commandName, args, cmd, bot, &trigger)
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
	rawCommand = strings.ToLower(parts[0])
	args = parts[1:]
	return
}
