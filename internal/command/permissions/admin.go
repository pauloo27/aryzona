package permissions

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

var MustBeAdmin = &command.Permission{
	Name: "be server admin",
	Checker: func(ctx *command.Context) bool {
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
