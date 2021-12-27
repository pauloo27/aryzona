package validations

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
)

var MustHaveVoicerOnGuiild = command.CommandValidation{
	Name: "have voicer on guild",
	Checker: func(ctx *command.CommandContext) (bool, string) {
		vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if vc == nil {
			return false, "Bot is not connect to a voice channel"
		}
		return true, ""
	},
}
