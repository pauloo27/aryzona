package listeners

import (
	"strings"

	"github.com/Pauloo27/aryzona/command"
	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Message.Content, command.Prefix) {
		return
	}

	rawCommand := strings.TrimPrefix(strings.Split(m.Content, " ")[0], command.Prefix)
	args := strings.Split(
		strings.TrimPrefix(strings.TrimPrefix(m.Content, command.Prefix+rawCommand), " "), " ",
	)
	command.HandleCommand(strings.ToLower(rawCommand), args, s, m)
}
