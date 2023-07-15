package validations

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
)

var MustHaveVoicerOnGuild = &command.Validation{
	Name: "mustHaveVoicerOnGuild",
	Checker: func(ctx *command.Context) (bool, string) {
		vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if vc == nil {
			return false, ctx.Lang.Validations.MustHaveVoicerOnGuild.BotNotConnected.Str()
		}
		ctx.Locals["vc"] = vc
		return true, ""
	},
}

var MustBePlaying = &command.Validation{
	Name:      "mustBePlaying",
	DependsOn: []*command.Validation{MustHaveVoicerOnGuild},
	Checker: func(ctx *command.Context) (bool, string) {
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

var MustBeOnVoiceChannel = &command.Validation{
	Name: "mustBeOnVoiceChannel",
	Checker: func(ctx *command.Context) (bool, string) {
		voiceState, err := ctx.Bot.FindUserVoiceState(ctx.GuildID, ctx.AuthorID)
		if err != nil {
			return false, ctx.Lang.Validations.MustBeOnVoiceChannel.YouAreNotInVoiceChannel.Str()
		}
		ctx.Locals["authorVoiceChannelID"] = voiceState.Channel().ID()
		return true, ""
	},
}

var MustBeOnAValidVoiceChannel = &command.Validation{
	Name:      "mustBeOnAValidVoiceChannel",
	DependsOn: []*command.Validation{MustBeOnVoiceChannel},
	Checker: func(ctx *command.Context) (bool, string) {
		vc, err := voicer.NewVoicerForUser(ctx.AuthorID, ctx.GuildID)
		if err != nil || !vc.CanConnect() {
			return false, ctx.Lang.Validations.MustBeOnAValidVoiceChannel.CannotConnectToChannel.Str()
		}
		ctx.Locals["vc"] = vc
		return true, ""
	},
}

var MustBeOnSameVoiceChannel = &command.Validation{
	Name:      "mustBeOnSameVoiceChannel",
	DependsOn: []*command.Validation{MustBeOnVoiceChannel, MustHaveVoicerOnGuild},
	Checker: func(ctx *command.Context) (bool, string) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		authorVoiceChannelID, found := ctx.Locals["authorVoiceChannelID"]
		if !found || *(vc.ChannelID) != authorVoiceChannelID.(string) {
			return false, ctx.Lang.Validations.MustBeOnSameVoiceChannel.NotInRightChannel.Str()
		}
		return true, ""
	},
}
