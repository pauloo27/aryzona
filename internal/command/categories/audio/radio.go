package audio

import (
	"errors"

	"github.com/Pauloo27/aryzona/internal/audio/dca"
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/providers/radio"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/logger"
)

var RadioCommand = command.Command{
	Name:        "radio",
	Description: "Plays a pre-defined radio",
	Deferred:    true,
	Parameters: []*command.CommandParameter{
		{
			Name:        "radio",
			Description: "radio name",
			Required:    false,
			Type:        parameters.ParameterString,
			ValidValuesFunc: func() []interface{} {
				ids := []interface{}{}
				for _, radio := range radio.GetRadioList() {
					ids = append(ids, radio.GetID())
				}
				return ids
			},
		},
	},
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 0 {
			listRadios(ctx, "Radio list:")
			return
		}
		if ok, msg := command.RunValidation(ctx, validations.MustBeOnAValidVoiceChannel); !ok {
			ctx.Error(msg)
			return
		}
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		radioID := ctx.Args[0].(string)
		channel := radio.GetRadioByID(radioID)

		if !vc.IsConnected() {
			if err := vc.Connect(); err != nil {
				ctx.Error("Cannot connect to your voice channel")
				return
			}
		}
		embed := buildPlayableInfoEmbed(channel, nil).WithTitle("Added to queue: " + channel.GetName())
		ctx.SuccessEmbed(embed)

		utils.Go(func() {
			if err := vc.AppendToQueue(channel); err != nil {
				if errors.Is(err, dca.ErrVoiceConnectionClosed) {
					return
				}
				ctx.Error(utils.Fmt("Cannot play stuff: %v", err))
				logger.Error(err)
				return
			}
		})
	},
}

func listRadios(ctx *command.CommandContext, title string) {
	embed := discord.NewEmbed().
		WithTitle(title)

	for _, channel := range radio.GetRadioList() {
		embed.WithFieldInline(channel.GetID(), channel.GetName())
	}

	embed.WithFooter(
		utils.Fmt(
			"Start a radio with '%s%s <name>' and '%sstop' when you are tired of it!",
			command.Prefix, ctx.UsedName, command.Prefix,
		),
	)

	ctx.SuccessEmbed(embed)
}
