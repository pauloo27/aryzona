package command

import (
	"github.com/bwmarrin/discordgo"
)

func HandleCommand(commandName string, s *discordgo.Session, m *discordgo.MessageCreate) {
	command, ok := commandMap[commandName]
	if !ok {
		return
	}

	command.Handler(&CommandContext{m.Message, m, s})
}
