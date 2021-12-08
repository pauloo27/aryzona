package dcgo

import (
	"errors"
	"time"

	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/event"
	"github.com/bwmarrin/discordgo"
)

func init() {
	discord.UseImplementation(&DcgoBot{
		d: &discordData{},
	})
}

type discordData struct {
	token     string
	s         *discordgo.Session
	startedAt *time.Time
	listeners []interface{}
}

type DcgoBot struct {
	d *discordData
}

func (b DcgoBot) Init(token string) error {
	b.d.token = token
	var err error
	b.d.s, err = discordgo.New("Bot " + b.d.token)
	return err
}

func (b DcgoBot) connect() error {
	now := time.Now()
	err := b.d.s.Open()
	if err != nil {
		return err
	}
	b.d.startedAt = &now
	return nil
}

func (b DcgoBot) Start() error {
	b.registerListeners()
	return b.connect()
}

func (b DcgoBot) StartedAt() *time.Time {
	return b.d.startedAt
}

func (b DcgoBot) disconnect() error {
	return b.d.s.Close()
}

func (b DcgoBot) Stop() error {
	return b.disconnect()
}

func (b DcgoBot) Self() (discord.User, error) {
	u := b.d.s.State.User
	return buildUser(u.ID), nil
}

func (b DcgoBot) Listen(eventType event.EventType, listener interface{}) error {
	var l interface{}
	switch eventType {
	case event.Ready:
		l = func(s *discordgo.Session, m *discordgo.Ready) {
			listener.(func(discord.BotAdapter))(b)
		}
	case event.MessageCreated:
		l = func(s *discordgo.Session, m *discordgo.MessageCreate) {
			msg := buildMessage(m.Message.ID, buildChannel(m.ChannelID, buildGuild(m.GuildID)), buildUser(m.Author.ID), m.Content)
			listener.(func(discord.BotAdapter, discord.Message))(b, msg)
		}
	default:
		return event.ErrEventNotSupported
	}
	b.d.listeners = append(b.d.listeners, l)
	return nil
}

func (b DcgoBot) registerListeners() {
	for _, listener := range b.d.listeners {
		b.d.s.AddHandler(listener)
	}
}

func (b DcgoBot) SendReplyMessage(m discord.Message, content string) (discord.Message, error) {
	msg, err := b.d.s.ChannelMessageSendReply(m.Channel().ID(), content, &discordgo.MessageReference{
		MessageID: m.ID(),
		ChannelID: m.Channel().ID(),
		GuildID:   m.Channel().Guild().ID(),
	})
	if err != nil {
		return nil, err
	}
	return buildMessage(msg.ID, buildChannel(msg.ChannelID, buildGuild(msg.GuildID)), buildUser(msg.Author.ID), msg.Content), nil
}

func (b DcgoBot) SendMessage(channelID string, message string) (discord.Message, error) {
	msg, err := b.d.s.ChannelMessageSend(channelID, message)
	if err != nil {
		return nil, err
	}
	return buildMessage(msg.ID, buildChannel(msg.ChannelID, buildGuild(msg.GuildID)), buildUser(msg.Author.ID), msg.Content), nil
}

func (b DcgoBot) SendReplyEmbedMessage(m discord.Message, embed *discord.Embed) (discord.Message, error) {
	msg, err := b.d.s.ChannelMessageSendComplex(m.Channel().ID(), &discordgo.MessageSend{
		Reference: &discordgo.MessageReference{
			MessageID: m.ID(),
			ChannelID: m.Channel().ID(),
			GuildID:   m.Channel().Guild().ID(),
		},
		Embed: buildEmbed(embed),
	})
	if err != nil {
		return nil, err
	}
	return buildMessage(msg.ID, buildChannel(msg.ChannelID, buildGuild(msg.GuildID)), buildUser(msg.Author.ID), msg.Content), nil
}

func (b DcgoBot) SendEmbedMessage(channelID string, embed *discord.Embed) (discord.Message, error) {
	msg, err := b.d.s.ChannelMessageSendEmbed(channelID, buildEmbed(embed))
	if err != nil {
		return nil, err
	}
	return buildMessage(msg.ID, buildChannel(msg.ChannelID, buildGuild(msg.GuildID)), buildUser(msg.Author.ID), msg.Content), nil
}

func (b DcgoBot) OpenChannelWithUser(userID string) (discord.Channel, error) {
	c, err := b.d.s.UserChannelCreate(userID)
	if err != nil {
		return nil, err
	}
	return buildChannel(c.ID, buildGuild(c.GuildID)), nil
}

func (b DcgoBot) Latency() time.Duration {
	return b.d.s.HeartbeatLatency()
}

func (b DcgoBot) OpenGuild(guildID string) (discord.Guild, error) {
	return nil, errors.New("not implemented yet")
}

func (b DcgoBot) JoinVoiceChannel(guildID, channelID string) (discord.VoiceChannel, error) {
	return nil, errors.New("not implemented yet")
}

func (b DcgoBot) FindUserVoiceState(guildID string, userID string) (discord.VoiceState, error) {
	return nil, errors.New("not implemented yet")
}

func (b DcgoBot) UpdatePresence(presence *discord.Presence) error {
	switch presence.Type {
	case discord.PresencePlaying:
		return b.d.s.UpdateGameStatus(0, presence.Title)
	case discord.PresenceListening:
		return b.d.s.UpdateListeningStatus(presence.Title)
	case discord.PresenceStreaming:
		return b.d.s.UpdateStreamingStatus(0, presence.Title, presence.Extra)
	default:
		return errors.New("invalid presence type")
	}
}
