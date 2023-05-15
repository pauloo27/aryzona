package arkw

import (
	"time"

	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/event"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/google/uuid"
)

func init() {
	eventCache := ttlcache.NewCache()
	_ = eventCache.SetTTL(1 * time.Minute)
	discord.UseImplementation(&ArkwBot{
		discordData: &discordData{
			prevChannelCache: eventCache,
		},
	})
}

type eventListener struct {
	handler  any
	preEvent bool
}

func (b ArkwBot) Listen(eventType event.EventType, listener any) error {
	var l any
	switch eventType {
	case event.Ready:
		l = func(*gateway.ReadyEvent) {
			listener.(func(discord.BotAdapter))(b)
		}
	case event.MessageCreated:
		b.indents = append(b.indents, gateway.IntentGuildMessages, gateway.IntentDirectMessages)
		l = func(m *gateway.MessageCreateEvent) {
			cType := model.ChannelTypeGuild
			if m.GuildID.String() == "" {
				cType = model.ChannelTypeDirect
			}
			msg := buildMessage(
				m.ID.String(),
				buildChannel(
					m.ChannelID.String(),
					buildGuild(m.GuildID.String()),
					cType,
				),
				buildUser(m.Author.ID.String()),
				m.Content,
			)
			listener.(func(discord.BotAdapter, model.Message))(b, msg)
		}
	case event.VoiceStateUpdated:
		eventID := uuid.New().String()
		b.indents = append(b.indents, gateway.IntentGuildVoiceStates, gateway.IntentGuilds)

		registerPreVoiceStateUpdateListener(b, eventID)

		l = func(m *gateway.VoiceStateUpdateEvent) {
			user := buildUser(m.UserID.String())
			var prevCh, curCh model.VoiceChannel
			if m.ChannelID.IsValid() {
				curCh = buildVoiceChannel(m.ChannelID.String(), buildGuild(m.GuildID.String()))
			}
			possiblePrevCh, err := b.prevChannelCache.Get(eventID)
			if err == nil {
				prevCh, _ = possiblePrevCh.(model.VoiceChannel)
			}

			listener.(func(discord.BotAdapter, model.User, model.VoiceChannel, model.VoiceChannel))(b, user, prevCh, curCh)
		}
	default:
		return event.ErrEventNotSupported
	}
	b.listeners = append(b.listeners, &eventListener{handler: l, preEvent: false})
	return nil
}

func (b ArkwBot) registerListeners() {
	for _, intent := range b.indents {
		b.s.AddIntents(intent)
	}
	b.s.PreHandler = handler.New()
	for _, l := range b.listeners {
		if l.preEvent {
			b.s.PreHandler.AddSyncHandler(l.handler)
		} else {
			b.s.AddHandler(l.handler)
		}
	}
}

func registerPreVoiceStateUpdateListener(b ArkwBot, eventID string) {
	// cache the state before the state update
	b.listeners = append(b.listeners, &eventListener{handler: func(m *gateway.VoiceStateUpdateEvent) {
		var prevCh model.VoiceChannel
		voiceState, err := b.FindUserVoiceState(m.GuildID.String(), m.UserID.String())
		if err == nil && voiceState.Channel().ID() != "" {
			cType := model.ChannelTypeGuild
			if voiceState.Channel().Guild().ID() == "" {
				cType = model.ChannelTypeDirect
			}
			prevCh = buildChannel(voiceState.Channel().ID(), buildGuild(voiceState.Channel().Guild().ID()), cType)
		}
		_ = b.prevChannelCache.Set(eventID, prevCh)
	}, preEvent: true})
}
