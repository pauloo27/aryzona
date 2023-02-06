package audio

import (
	"errors"

	"github.com/Pauloo27/aryzona/internal/audio/dca"
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/core/routine"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/aryzona/internal/providers/radio"
	"github.com/Pauloo27/logger"
)

var RadioCommand = command.Command{
	Name:     "radio",
	Deferred: true,
	Parameters: []*command.CommandParameter{
		{
			Name:     "radio",
			Required: false,
			Type:     parameters.ParameterString,
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
		t := ctx.T.(*i18n.CommandRadio)

		if len(ctx.Args) == 0 {
			listRadios(ctx, t.ListTitle.Str())
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
				ctx.Error(t.CannotConnect.Str())
				logger.Error(err)
				return
			}
		} else {
			authorVoiceChannelID, found := ctx.Locals["authorVoiceChannelID"]
			if !found || *(vc.ChannelID) != authorVoiceChannelID.(string) {
				ctx.Error(t.NotInRightChannel.Str())
				return
			}
		}
		embed := buildPlayableInfoEmbed(channel, nil, ctx.AuthorID, t.PlayingInfo).WithTitle(t.AddedToQueue.Str(channel.GetName()))
		ctx.SuccessEmbed(embed)

		routine.Go(func() {
			if err := vc.AppendToQueue(ctx.AuthorID, channel); err != nil {
				if errors.Is(err, dca.ErrVoiceConnectionClosed) {
					return
				}
				ctx.Errorf(t.SomethingWentWrong.Str())
				logger.Error(err)
				return
			}
		})
	},
}

func listRadios(ctx *command.CommandContext, title string) {
	t := ctx.T.(*i18n.CommandRadio)

	embed := model.NewEmbed().
		WithTitle(title)

	for _, channel := range radio.GetRadioList() {
		embed.WithFieldInline(channel.GetID(), channel.GetName())
	}

	embed.WithFooter(
		t.ListFooter.Str(
			command.Prefix, ctx.UsedName,
		),
	)

	ctx.SuccessEmbed(embed)
}
