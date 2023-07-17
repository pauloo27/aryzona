package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/radio"
	"github.com/pauloo27/logger"
)

var RadioCommand = command.Command{
	Name:     "radio",
	Deferred: true,
	Parameters: []*command.Parameter{
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
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandRadio)

		if len(ctx.Args) == 0 {
			return listRadios(ctx, t.ListTitle.Str())
		}

		if ok, msg := command.RunValidation(ctx, validations.MustBeOnAValidVoiceChannel); !ok {
			return ctx.Error(msg)
		}
		vc := ctx.Locals["vc"].(*voicer.Voicer)

		radioID := ctx.Args[0].(string)
		channel := radio.GetRadioByID(radioID)

		if !vc.IsConnected() {
			if err := vc.Connect(); err != nil {
				logger.Error(err)
				return ctx.Error(t.CannotConnect.Str())
			}
		} else {
			ok, msg := validations.MustBeOnSameVoiceChannel.Checker(ctx)
			if !ok {
				logger.Error(msg)
				return ctx.Error(msg)
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

		vc.AppendToQueue(ctx.AuthorID, channel)

		return ctx.SuccessEmbed(embed)
	},
}

func listRadios(ctx *command.Context, title string) command.Result {
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

	return ctx.SuccessEmbed(embed)
}
