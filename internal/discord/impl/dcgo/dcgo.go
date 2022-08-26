package dcgo

import (
	"errors"
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/event"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/bwmarrin/discordgo"
)

func init() {
	discord.UseImplementation(&DcgoBot{
		d: &discordData{},
	})
}

type discordData struct {
	listeners []interface{}
	token     string
	s         *discordgo.Session
	startedAt *time.Time
}

type DcgoBot struct {
	d *discordData
}

func (b DcgoBot) Implementation() string {
	return "Discordgo"
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

func (b DcgoBot) CountUsersInVoiceChannel(ch model.VoiceChannel) (count int) {
	g, err := b.d.s.State.Guild(ch.Guild().ID())
	if err != nil {
		return 0
	}
	for _, voice := range g.VoiceStates {
		if voice.ChannelID == ch.ID() {
			count++
		}
	}
	return
}

func (b DcgoBot) disconnect() error {
	return b.d.s.Close()
}

func (b DcgoBot) Stop() error {
	return b.disconnect()
}

func (b DcgoBot) Self() (model.User, error) {
	u := b.d.s.State.User
	return buildUser(u.ID), nil
}

func (b DcgoBot) GuildCount() int {
	return len(b.d.s.State.Guilds)
}

func (b DcgoBot) RegisterSlashCommands() error {
	return registerCommands(b)
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
			cType := model.ChannelTypeGuild
			if m.GuildID == "" {
				cType = model.ChannelTypeDirect
			}
			msg := buildMessage(m.Message.ID, buildChannel(m.ChannelID, buildGuild(m.GuildID), cType), buildUser(m.Author.ID), m.Content)
			listener.(func(discord.BotAdapter, model.Message))(b, msg)
		}
	case event.VoiceStateUpdated:
		l = func(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
			var prevCh, curCh model.VoiceChannel
			user := buildUser(e.UserID)
			if e.BeforeUpdate != nil {
				prevCh = buildVoiceChannel(e.BeforeUpdate.ChannelID, buildGuild(e.BeforeUpdate.GuildID))
			}
			if e.ChannelID != "" {
				curCh = buildVoiceChannel(e.ChannelID, buildGuild(e.GuildID))
			}
			listener.(func(discord.BotAdapter, model.User, model.VoiceChannel, model.VoiceChannel))(b, user, prevCh, curCh)
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

func (b DcgoBot) SendReplyMessage(m model.Message, content string) (model.Message, error) {
	msg, err := b.d.s.ChannelMessageSendReply(m.Channel().ID(), content, &discordgo.MessageReference{
		MessageID: m.ID(),
		ChannelID: m.Channel().ID(),
		GuildID:   m.Channel().Guild().ID(),
	})
	if err != nil {
		return nil, err
	}
	cType := model.ChannelTypeGuild
	if msg.GuildID == "" {
		cType = model.ChannelTypeDirect
	}
	return buildMessage(msg.ID, buildChannel(msg.ChannelID, buildGuild(msg.GuildID), cType), buildUser(msg.Author.ID), msg.Content), nil
}

func (b DcgoBot) SendMessage(channelID string, message string) (model.Message, error) {
	msg, err := b.d.s.ChannelMessageSend(channelID, message)
	if err != nil {
		return nil, err
	}
	cType := model.ChannelTypeGuild
	if msg.GuildID == "" {
		cType = model.ChannelTypeDirect
	}
	return buildMessage(msg.ID, buildChannel(msg.ChannelID, buildGuild(msg.GuildID), cType), buildUser(msg.Author.ID), msg.Content), nil
}

func (b DcgoBot) SendReplyEmbedMessage(m model.Message, embed *discord.Embed) (model.Message, error) {
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
	cType := model.ChannelTypeGuild
	if msg.GuildID == "" {
		cType = model.ChannelTypeDirect
	}
	return buildMessage(msg.ID, buildChannel(msg.ChannelID, buildGuild(msg.GuildID), cType), buildUser(msg.Author.ID), msg.Content), nil
}

func (b DcgoBot) SendEmbedMessage(channelID string, embed *discord.Embed) (model.Message, error) {
	msg, err := b.d.s.ChannelMessageSendEmbed(channelID, buildEmbed(embed))
	if err != nil {
		return nil, err
	}
	cType := model.ChannelTypeGuild
	if msg.GuildID == "" {
		cType = model.ChannelTypeDirect
	}
	return buildMessage(msg.ID, buildChannel(msg.ChannelID, buildGuild(msg.GuildID), cType), buildUser(msg.Author.ID), msg.Content), nil
}

func (b DcgoBot) OpenChannelWithUser(userID string) (model.TextChannel, error) {
	c, err := b.d.s.UserChannelCreate(userID)
	if err != nil {
		return nil, err
	}
	return buildChannel(c.ID, buildGuild(c.GuildID), model.ChannelTypeDirect), nil
}

func (b DcgoBot) Latency() time.Duration {
	return b.d.s.HeartbeatLatency()
}

func (b DcgoBot) OpenGuild(guildID string) (model.Guild, error) {
	g, err := b.d.s.Guild(guildID)
	if err != nil {
		return nil, err
	}
	return buildGuild(g.ID), nil
}

func (b DcgoBot) JoinVoiceChannel(guildID, channelID string) (model.VoiceConnection, error) {
	vc, err := b.d.s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}
	return buildVoiceConnection(vc), nil
}

func (b DcgoBot) FindUserVoiceState(guildID, userID string) (model.VoiceState, error) {
	state, err := b.d.s.State.VoiceState(guildID, userID)
	if err != nil {
		return nil, err
	}
	return buildVoiceState(buildVoiceChannel(state.ChannelID, buildGuild(guildID))), nil
}

func (b DcgoBot) GetMember(guilID, userID string) (model.Member, error) {
	return nil, errors.New("not implemented yet")
}

func (b DcgoBot) UpdatePresence(presence *model.Presence) error {
	switch presence.Type {
	case model.PresencePlaying:
		return b.d.s.UpdateGameStatus(0, presence.Title)
	case model.PresenceListening:
		return b.d.s.UpdateListeningStatus(presence.Title)
	case model.PresenceStreaming:
		return b.d.s.UpdateStreamingStatus(0, presence.Title, presence.Extra)
	default:
		return errors.New("invalid presence type")
	}
}
