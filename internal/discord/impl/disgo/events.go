package disgo

import (
	"fmt"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/event"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

func (d *DisgoBot) Listen(eventType event.EventType, handlerFunc any) error {
	switch eventType {
	case event.Ready:
		wrapHandler := func(e *events.Ready) {
			handlerFunc.(func(discord.BotAdapter))(d)
		}
		d.opts = append(d.opts, bot.WithEventListenerFunc(wrapHandler))
	case event.MessageCreated:
		d.intents = append(
			d.intents,
			gateway.IntentGuildMessages, gateway.IntentDirectMessages, gateway.IntentMessageContent,
		)
		wrapHandler := func(e *events.MessageCreate) {
			channelType := model.ChannelTypeDirect
			guildID := ""
			if e.Message.GuildID != nil {
				channelType = model.ChannelTypeGuild
				guildID = e.Message.GuildID.String()
			}

			msg := buildMessage(
				e.Message.ID.String(),
				buildChannel(
					e.ChannelID.String(),
					buildGuild(guildID),
					channelType,
				),
				buildUser(e.Message.Author.ID.String()),
				e.Message.Content,
			)

			handlerFunc.(func(discord.BotAdapter, model.Message))(d, msg)
		}
		d.opts = append(d.opts, bot.WithEventListenerFunc(wrapHandler))
	case event.VoiceStateUpdated: // TODO:
	default:
		return fmt.Errorf("unexpected event type %d", eventType)
	}

	return nil
}
