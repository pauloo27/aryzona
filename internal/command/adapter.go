package command

import (
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"
)

type Adapter struct {
	GuildID, AuthorID string
	Reply             func(*CommandContext, string) error
	DeferResponse     func() error
	ReplyEmbed        func(*CommandContext, *discord.Embed) error
	Member            model.Member
}
