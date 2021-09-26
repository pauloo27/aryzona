package listeners

import (
	"strings"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.ID == "214486492909666305" {
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			logger.Error(err)
		}
	}

	if !strings.HasPrefix(m.Message.Content, command.Prefix) {
		return
	}

	rawCommand := strings.TrimPrefix(strings.Split(m.Content, " ")[0], command.Prefix)
	args := strings.Split(
		strings.TrimPrefix(strings.TrimPrefix(m.Content, command.Prefix+rawCommand), " "), " ",
	)
	if len(args) == 1 && args[0] == "" {
		args = []string{}
	}
	event := command.Event{
		Message: m.Message,
	}
	command.HandleCommand(strings.ToLower(rawCommand), args, s, &event)
}
