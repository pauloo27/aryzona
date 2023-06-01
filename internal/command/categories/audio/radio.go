package audio

import (
	"errors"

	"github.com/pauloo27/aryzona/internal/audio/dca"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/core/routine"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/radio"
	"github.com/pauloo27/logger"
)

var RadioCommand = command.Command{
	Name:     "radio",
	Deferred: true,
	Parameters: []*command.CommandParameter{
		{
			Name:     "radio",
			Required: false,
			Type:     parameters.ParameterString,
			ValidValuesFunc: func() []any {
				ids := []any{}
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
			ok, msg := validations.MustBeOnSameVoiceChannel.Checker(ctx)
			if !ok {
				logger.Error(msg)
				return
			}
		}
		embed := buildPlayableInfoEmbed(
			PlayableInfo{
				Playable:    channel,
				RequesterID: ctx.AuthorID,
				T:           t.PlayingInfo,
				Common:      t.Common,
			},
		).WithTitle(t.AddedToQueue.Str(channel.GetName()))
		ctx.SuccessEmbed(embed)

		routine.GoAndRecover(func() {
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
