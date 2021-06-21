package audio

import (
	"github.com/Pauloo27/aryzona/command"
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
		ctx.Success("sorry, but right now I can only dance")
		// TODO: safe connect to channel
		// TODO: play audio: Probably a common audio stuff providers dont mess up
		// TODO: leave when empty
		// TODO: stop command: even more commands?
	},
}
