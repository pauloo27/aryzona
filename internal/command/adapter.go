package command

import (
	"github.com/pauloo27/aryzona/internal/discord/model"
)

type Adapter struct {
	GuildID, AuthorID string
	Reply             func(*CommandContext, string) error
	Edit              func(*CommandContext, string) error
	DeferResponse     func() error
	ReplyEmbed        func(*CommandContext, *model.Embed) error
	ReplyComplex      func(*CommandContext, *model.ComplexMessage) error
	EditEmbed         func(*CommandContext, *model.Embed) error
	EditComplex       func(*CommandContext, *model.ComplexMessage) error
}
