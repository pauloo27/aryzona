package validations

import (
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
)

var MustHaveVoicerOnGuild = &command.CommandValidation{
	Name: "mustHaveVoicerOnGuild",
	Checker: func(ctx *command.CommandContext) (bool, string) {
		vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if vc == nil {
			return false, ctx.Lang.Validations.MustHaveVoicerOnGuild.BotNotConnected.Str()
		}
		ctx.Locals["vc"] = vc
		return true, ""
	},
}

var MustBePlaying = &command.CommandValidation{
	Name:      "mustBePlaying",
	DependsOn: []*command.CommandValidation{MustHaveVoicerOnGuild},
	Checker: func(ctx *command.CommandContext) (bool, string) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		item := vc.Playing()
		if item == nil {
			return false, ctx.Lang.Validations.MustBePlaying.NothingPlaying.Str()
		}
		ctx.Locals["playing"] = item.Playable
		ctx.Locals["requesterID"] = item.Requester

		return true, ""
	},
}

var MustBeOnVoiceChannel = &command.CommandValidation{
	Name: "mustBeOnVoiceChannel",
	Checker: func(ctx *command.CommandContext) (bool, string) {
		voiceState, err := ctx.Bot.FindUserVoiceState(ctx.GuildID, ctx.AuthorID)
		if err != nil {
			return false, ctx.Lang.Validations.MustBeOnVoiceChannel.YouAreNotInVoiceChannel.Str()
		}
		ctx.Locals["authorVoiceChannelID"] = voiceState.Channel().ID()
		return true, ""
	},
}

var MustBeOnAValidVoiceChannel = &command.CommandValidation{
	Name:      "mustBeOnAValidVoiceChannel",
	DependsOn: []*command.CommandValidation{MustBeOnVoiceChannel},
	Checker: func(ctx *command.CommandContext) (bool, string) {
		vc, err := voicer.NewVoicerForUser(ctx.AuthorID, ctx.GuildID)
		if err != nil || !vc.CanConnect() {
			return false, ctx.Lang.Validations.MustBeOnAValidVoiceChannel.CannotConnectToChannel.Str()
		}
		ctx.Locals["vc"] = vc
		return true, ""
	},
}

var MustBeOnSameVoiceChannel = &command.CommandValidation{
	Name:      "mustBeOnSameVoiceChannel",
	DependsOn: []*command.CommandValidation{MustBeOnVoiceChannel, MustHaveVoicerOnGuild},
	Checker: func(ctx *command.CommandContext) (bool, string) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		authorVoiceChannelID, found := ctx.Locals["authorVoiceChannelID"]
		if !found || *(vc.ChannelID) != authorVoiceChannelID.(string) {
			return false, ctx.Lang.Validations.MustBeOnSameVoiceChannel.NotInRightChannel.Str()
		}
		return true, ""
	},
}
