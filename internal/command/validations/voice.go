package validations

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
)

var MustHaveVoicerOnGuild = &command.CommandValidation{
	Name: "have voicer on guild",
	Checker: func(ctx *command.CommandContext) (bool, string) {
		vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if vc == nil {
			return false, "Bot is not connect to a voice channel"
		}
		ctx.Locals["vc"] = vc
		return true, ""
	},
}

var MustBePlaying = &command.CommandValidation{
	Name:      "be playing something on voicer",
	DependsOn: []*command.CommandValidation{MustHaveVoicerOnGuild},
	Checker: func(ctx *command.CommandContext) (bool, string) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := vc.Playing()

		if playing == nil {
			return false, "Nothing playing"
		}
		ctx.Locals["playing"] = playing

		return true, ""
	},
}
