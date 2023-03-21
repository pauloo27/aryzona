package permissions

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/model"
)

var MustBeAdmin = command.CommandPermission{
	Name: "be server admin",
	Checker: func(ctx *command.CommandContext) bool {
		if ctx.GuildID == "" {
			return false
		}

		member, err := ctx.Bot.GetMember(ctx.GuildID, "", ctx.AuthorID)
		if err != nil {
			return false
		}

		return member.Permissions().Has(model.PermissionAdministrator)
	},
}
