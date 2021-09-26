package command

import "github.com/bwmarrin/discordgo"

type Event struct {
	Reply             func(string) error
	ReplyEmbed        func(*discordgo.MessageEmbed) error
	GuildID, AuthorID string
}
