package tools

import (
	"github.com/google/uuid"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

var UUIDCommand = command.Command{
	Name: "uuid",
	Handler: func(ctx *command.Context) {
		id := uuid.New()
		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle("UUID v4").
				WithDescription(id.String()),
		)
	},
}
