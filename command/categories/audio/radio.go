package audio

import (
	"github.com/Pauloo27/aryzona/audio/dca"
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

	embed.Footer(
		utils.Fmt(
			"Use `%sradio <name>` and `%sradio stop` when you are tired of it!",
			command.Prefix, command.Prefix,
		), "",
	)

	ctx.SuccesEmbed(embed.Build())
}

var RadioCommand = command.Command{
	Name:        "radio",
	Description: "Plays a pre-defined radio",
	Arguments: []*command.CommandArgument{
		{
			Name:     "radio name",
			Required: false,
			Type:     command.ArgumentString,
			ValidValuesFunc: func() []interface{} {
				ids := []interface{}{}
				for _, radio := range radio.GetRadioList() {
					ids = append(ids, radio.Id)
				}
				return append(ids, "stop")
			},
		},
	},
	Permission: &permissions.BeOwner,
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 0 {
			listRadios(ctx, "Radio list:")
			return
		}

		vc, err := voicer.NewVoicerForUser(ctx.Message.Author.ID, ctx.Message.GuildID)
		if err != nil {
			ctx.Error("Cannot create voicer")
			return
		}

		var channel *radio.RadioChannel
		radioId := ctx.Args[0].(string)

		if radioId == "stop" {
			if !vc.IsConnected() || !vc.IsPlaying() {
				ctx.Error("Already stopped")
			} else {
				err = vc.Disconnect()
				if err != nil {
					ctx.Error(utils.Fmt("Cannot disconnect: %v", err))
				} else {
					ctx.Success("Disconnected")
				}
			}
			return
		} else {
			channel = radio.GetRadioById(radioId)
		}

		if !vc.CanConnect() {
			ctx.Error("You are not in a voice channel")
			return
		}
		if err = vc.Connect(); err != nil {
			ctx.Error("Cannot connect to your voice channel")
			return
		}
		go func() {
			if err = vc.Play(channel); err != nil {
				if is, vErr := utils.IsErrore(err); is {
					if vErr.ID == dca.ERR_VOICE_CONNECTION_CLOSED.ID {
						return
					}
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
