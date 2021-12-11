package arkw

import (
	"context"
	"errors"
	"time"

	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/event"
	dc "github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

func init() {
	discord.UseImplementation(&ArkwBot{
		d: &discordData{},
	})
}

type discordData struct {
	token     string
	startedAt *time.Time
	listeners []interface{}
	indents   []gateway.Intents
	s         *session.Session
}

type ArkwBot struct {
	d *discordData
}

func (b ArkwBot) Init(token string) error {
	b.d.token = token
	var err error
	b.d.s, err = session.New("Bot " + b.d.token)
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
	return 0
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
	v, _ := b.d.s.Guilds(0)
	return len(v)
}

func (b ArkwBot) RegisterSlashCommands() error {
	return errors.New("not implemented yet")
}

func (b ArkwBot) Listen(eventType event.EventType, listener interface{}) error {
	var l interface{}
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
		b.d.indents = append(b.d.indents, gateway.IntentGuildVoiceStates)
		return nil // FIXME
	default:
		return event.ErrEventNotSupported
	}
	b.d.listeners = append(b.d.listeners, l)
	return nil
}

func (b ArkwBot) registerListeners() {
	for _, intent := range b.d.indents {
		b.d.s.AddIntents(intent)
	}
	for _, l := range b.d.listeners {
		b.d.s.AddHandler(l)
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
	return buildChannel(dm.ID.String(), buildGuild("")), nil
}

func (b ArkwBot) Latency() time.Duration {
	return 0
}

func (b ArkwBot) OpenGuild(guildID string) (discord.Guild, error) {
	return nil, errors.New("not implemented yet")
}

func (b ArkwBot) JoinVoiceChannel(guildID, channelID string) (discord.VoiceConnection, error) {
	return nil, errors.New("not implemented yet")
}

func (b ArkwBot) FindUserVoiceState(guildID, userID string) (discord.VoiceState, error) {
	return nil, errors.New("not implemented yet")
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
