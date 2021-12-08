package command

import (
	"github.com/Pauloo27/aryzona/discord"
)

type Event struct {
	Reply             func(string) error
	ReplyEmbed        func(*discord.Embed) error
	GuildID, AuthorID string
}
