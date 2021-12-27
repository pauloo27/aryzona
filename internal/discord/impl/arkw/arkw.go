package arkw

import (
	"context"
	"errors"
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/event"
	"github.com/ReneKroon/ttlcache/v2"
	dc "github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/google/uuid"
)

func init() {
	cache := ttlcache.NewCache()
	_ = cache.SetTTL(1 * time.Minute)
	discord.UseImplementation(&ArkwBot{
		d: &discordData{
			prevChannelCache: cache,
		},
	})
}

type eventListener struct {
	preHandler bool
	handler    interface{}
}

type discordData struct {
	token            string
	startedAt        *time.Time
	listeners        []*eventListener
	indents          []gateway.Intents
	s                *state.State
	prevChannelCache *ttlcache.Cache
}

type ArkwBot struct {
	d *discordData
}

func (b ArkwBot) Implementation() string {
	return "Arikawa"
}

func (b ArkwBot) Init(token string) error {
	b.d.token = token
	var err error
	b.d.s, err = state.New("Bot " + b.d.token)
	return err
}

func (b ArkwBot) connect() error {
	now := time.Now()
	err := b.d.s.Open(context.Background())
	if err != nil {
		return err
	}
	b.d.startedAt = &now
	return nil
}

func (b ArkwBot) Start() error {
	b.registerListeners()
	return b.connect()
}

func (b ArkwBot) StartedAt() *time.Time {
	return b.d.startedAt
}

func (b ArkwBot) CountUsersInVoiceChannel(ch discord.VoiceChannel) (count int) {
	sf, err := dc.ParseSnowflake(ch.Guild().ID())
	if err != nil {
		return
	}
	states, err := b.d.s.VoiceStates(dc.GuildID(sf))
	if err != nil {
		return
	}
	for _, state := range states {
		if state.ChannelID.String() == ch.ID() {
			count++
		}
	}
	return
}

func (b ArkwBot) disconnect() error {
	return b.d.s.Close()
}

func (b ArkwBot) Stop() error {
	return b.disconnect()
}

func (b ArkwBot) Self() (discord.User, error) {
	u, err := b.d.s.Me()
	if err != nil {
		return nil, err
	}
	return buildUser(u.ID.String()), nil
}

func (b ArkwBot) GuildCount() int {
	v, _ := b.d.s.Guilds()
	return len(v)
}

func (b ArkwBot) RegisterSlashCommands() error {
	return registerCommands(b)
}

func (b ArkwBot) Listen(eventType event.EventType, listener interface{}) error {
	var l interface{}
	pre := false
	switch eventType {
	case event.Ready:
		l = func(*gateway.ReadyEvent) {
			listener.(func(discord.BotAdapter))(b)
		}
	case event.MessageCreated:
		b.d.indents = append(b.d.indents, gateway.IntentGuildMessages)
		b.d.indents = append(b.d.indents, gateway.IntentDirectMessages)
		l = func(m *gateway.MessageCreateEvent) {
			msg := buildMessage(m.ID.String(), buildChannel(m.ChannelID.String(), buildGuild(m.GuildID.String())), buildUser(m.Author.ID.String()), m.Content)
			listener.(func(discord.BotAdapter, discord.Message))(b, msg)
		}
	case event.VoiceStateUpdated:
		eventID := uuid.New().String()
		// add helper listener
		b.d.listeners = append(b.d.listeners, &eventListener{handler: func(m *gateway.VoiceStateUpdateEvent) {
			var prevCh discord.VoiceChannel
			voiceState, err := b.FindUserVoiceState(m.GuildID.String(), m.UserID.String())
			if err == nil && voiceState.Channel().ID() != "" {
				prevCh = buildChannel(voiceState.Channel().ID(), buildGuild(voiceState.Channel().Guild().ID()))
			}
			_ = b.d.prevChannelCache.Set(eventID, prevCh)
		}, preHandler: true})

		pre = false
		b.d.indents = append(b.d.indents, gateway.IntentGuildVoiceStates)
		b.d.indents = append(b.d.indents, gateway.IntentGuilds)

		l = func(m *gateway.VoiceStateUpdateEvent) {
			user := buildUser(m.UserID.String())
			var prevCh, curCh discord.VoiceChannel
			if m.ChannelID.IsValid() {
				curCh = buildVoiceChannel(m.ChannelID.String(), buildGuild(m.GuildID.String()))
			}
			possiblePrevCh, err := b.d.prevChannelCache.Get(eventID)
			if err == nil {
				prevCh, _ = possiblePrevCh.(discord.VoiceChannel)
			}

			listener.(func(discord.BotAdapter, discord.User, discord.VoiceChannel, discord.VoiceChannel))(b, user, prevCh, curCh)
		}
	default:
		return event.ErrEventNotSupported
	}
	b.d.listeners = append(b.d.listeners, &eventListener{handler: l, preHandler: pre})
	return nil
}

func (b ArkwBot) registerListeners() {
	for _, intent := range b.d.indents {
		b.d.s.AddIntents(intent)
	}
	b.d.s.PreHandler = handler.New()
	for _, l := range b.d.listeners {
		if l.preHandler {
			b.d.s.PreHandler.AddSyncHandler(l.handler)
		} else {
			b.d.s.AddHandler(l.handler)
		}
	}
}

func (b ArkwBot) SendReplyMessage(m discord.Message, content string) (discord.Message, error) {
	sf, err := dc.ParseSnowflake(m.Channel().ID())
	if err != nil {
		return nil, err
	}
	msg, err := b.d.s.SendMessage(dc.ChannelID(sf), content)
	if err != nil {
		return nil, err
	}
	return buildMessage(
		msg.ID.String(), buildChannel(m.Channel().ID(), buildGuild(msg.GuildID.String())),
		buildUser(msg.Author.ID.String()),
		content,
	), nil
}

func (b ArkwBot) SendMessage(channelID string, message string) (discord.Message, error) {
	sf, err := dc.ParseSnowflake(channelID)
	if err != nil {
		return nil, err
	}
	msg, err := b.d.s.SendMessage(dc.ChannelID(sf), message)
	if err != nil {
		return nil, err
	}
	return buildMessage(
		msg.ID.String(), buildChannel(channelID, buildGuild(msg.GuildID.String())),
		buildUser(msg.Author.ID.String()),
		message,
	), nil
}

func (b ArkwBot) SendReplyEmbedMessage(m discord.Message, embed *discord.Embed) (discord.Message, error) {
	chSf, err := dc.ParseSnowflake(m.Channel().ID())
	if err != nil {
		return nil, err
	}
	refSf, err := dc.ParseSnowflake(m.ID())
	if err != nil {
		return nil, err
	}
	msg, err := b.d.s.SendEmbedReply(dc.ChannelID(chSf), dc.MessageID(refSf), buildEmbed(embed))
	if err != nil {
		return nil, err
	}
	return buildMessage(
		msg.ID.String(), buildChannel(m.Channel().ID(), buildGuild(msg.GuildID.String())),
		buildUser(msg.Author.ID.String()),
		msg.Content,
	), nil
}

func (b ArkwBot) SendEmbedMessage(channelID string, embed *discord.Embed) (discord.Message, error) {
	sf, err := dc.ParseSnowflake(channelID)
	if err != nil {
		return nil, err
	}
	msg, err := b.d.s.SendEmbeds(dc.ChannelID(sf), buildEmbed(embed))
	if err != nil {
		return nil, err
	}
	return buildMessage(
		msg.ID.String(), buildChannel(channelID, buildGuild(msg.GuildID.String())),
		buildUser(msg.Author.ID.String()),
		msg.Content,
	), nil
}

func (b ArkwBot) OpenChannelWithUser(userID string) (discord.Channel, error) {
	sf, err := dc.ParseSnowflake(userID)
	if err != nil {
		return nil, err
	}
	dm, err := b.d.s.CreatePrivateChannel(dc.UserID(sf))
	if err != nil {
		return nil, err
	}
	return buildChannel(dm.ID.String(), buildGuild("")), nil
}

func (b ArkwBot) Latency() time.Duration {
	return b.d.s.PacerLoop.EchoBeat.Time().Sub(b.d.s.PacerLoop.SentBeat.Time())
}

func (b ArkwBot) OpenGuild(guildID string) (discord.Guild, error) {
	sf, err := dc.ParseSnowflake(guildID)
	if err != nil {
		return nil, err
	}
	g, err := b.d.s.Guild(dc.GuildID(sf))
	if err != nil {
		return nil, err
	}
	return buildGuild(g.ID.String()), nil
}

func (b ArkwBot) JoinVoiceChannel(guildID, channelID string) (discord.VoiceConnection, error) {
	vs, err := voice.NewSession(b.d.s)
	if err != nil {
		return nil, err
	}
	guildSf, err := dc.ParseSnowflake(guildID)
	if err != nil {
		return nil, err
	}
	channelSf, err := dc.ParseSnowflake(channelID)
	if err != nil {
		return nil, err
	}
	err = vs.JoinChannel(dc.GuildID(guildSf), dc.ChannelID(channelSf), false, false)
	if err != nil {
		return nil, err
	}
	return buildVoiceConnection(vs), nil
}

func (b ArkwBot) FindUserVoiceState(guildID, userID string) (discord.VoiceState, error) {
	guildSf, err := dc.ParseSnowflake(guildID)
	if err != nil {
		return nil, err
	}
	userSf, err := dc.ParseSnowflake(userID)
	if err != nil {
		return nil, err
	}
	vs, err := b.d.s.VoiceState(dc.GuildID(guildSf), dc.UserID(userSf))
	if err != nil {
		return nil, err
	}
	return buildVoiceState(buildVoiceChannel(vs.ChannelID.String(), buildGuild(guildID))), nil
}

func (b ArkwBot) UpdatePresence(presence *discord.Presence) error {
	var ty dc.ActivityType
	switch presence.Type {
	case discord.PresencePlaying:
		ty = dc.GameActivity
	case discord.PresenceListening:
		ty = dc.ListeningActivity
	case discord.PresenceStreaming:
		ty = dc.StreamingActivity
	default:
		return errors.New("invalid presence type")
	}
	return b.d.s.UpdateStatus(gateway.UpdateStatusData{
		Activities: []dc.Activity{
			{Name: presence.Title, URL: presence.Extra, Type: ty},
		},
	})
}
