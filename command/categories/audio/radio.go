package audio

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/permissions"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/logger"
	"github.com/Pauloo27/aryzona/providers/radio"
	"github.com/Pauloo27/aryzona/utils"
)

func listRadios(ctx *command.CommandContext, title string) {
	embed := utils.NewEmbedBuilder().
		Title(title)

	for _, channel := range radio.GetRadioList() {
		embed.Field(channel.Id, channel.Name)
	}

	embed.Footer("Use !radio <name> to listen to one!", "")

	ctx.SuccesEmbed(embed.Build())
}

var RadioCommand = command.Command{
	Name:        "radio",
	Description: "Plays a pre-defined radio",
	Permission:  &permissions.BeOwner,
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 0 {
			listRadios(ctx, "Radio list:")
			return
		}
		radioId := ctx.Args[0]
		channel := radio.GetRadioById(radioId)
		if channel == nil {
			listRadios(ctx, "Invalid radio id. Here are some valid ones:")
			return
		}
		vc, err := voicer.NewVoicerForUser(ctx.Message.Author.ID, ctx.Message.GuildID)
		if err != nil {
			ctx.Error("Cannot create voicer")
			return
		}
		if !vc.CanConnect() {
			ctx.Error("You are not in a voice channel")
			return
		}
		if err = vc.Connect(); err != nil {
			ctx.Error("Cannot  to your voice channel")
			return
		}
		go func() {
			if err = vc.Play(channel); err != nil {
				if is, vErr := voicer.IsVoicerError(err); is {
					ctx.Error(vErr.Message)
				} else {
					ctx.Error("Cannot play stuff")
					logger.Error(err.Error())
				}
				return
			}
		}()
	},
}
