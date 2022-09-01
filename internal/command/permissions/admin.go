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
		for _, r := range ctx.Member.Roles() {
			if r.Permissions().Has(model.PermissionAdministrator) {
				return true
			}
		}
		return false
	},
}
