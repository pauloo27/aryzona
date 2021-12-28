package validations

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
)

var MustHaveVoicerOnGuild = &command.CommandValidation{
	Description: "have voicer on guild",
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
	Description: "be playing something on voicer",
	DependsOn:   []*command.CommandValidation{MustHaveVoicerOnGuild},
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

var MustBeOnVoiceChannel = &command.CommandValidation{
	Description: "be connected to a voice channel",
	Checker: func(ctx *command.CommandContext) (bool, string) {
		if _, err := ctx.Bot.FindUserVoiceState(ctx.GuildID, ctx.AuthorID); err != nil {
			return false, "You are not in a voice channel"
		}
		return true, ""
	},
}

var MustBeOnAValidVoiceChannel = &command.CommandValidation{
	Description: "be connected to a valid voice channel",
	DependsOn:   []*command.CommandValidation{MustBeOnVoiceChannel},
	Checker: func(ctx *command.CommandContext) (bool, string) {
		vc, err := voicer.NewVoicerForUser(ctx.AuthorID, ctx.GuildID)
		if err != nil {
			return false, "Cannot create voicer"
		}
		if !vc.CanConnect() {
			return false, "Cannot connect to your voice channel"
		}
		ctx.Locals["vc"] = vc
		return true, ""
	},
}
