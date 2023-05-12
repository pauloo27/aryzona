package arkw

import (
	"context"
	"errors"
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/diamondburned/arikawa/v3/api"
	dc "github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/voice"
)

var _ discord.BotAdapter = ArkwBot{}

type discordData struct {
	listeners        []*eventListener
	indents          []gateway.Intents
	token            string
	startedAt        *time.Time
	s                *state.State
	prevChannelCache *ttlcache.Cache
}

type ArkwBot struct {
	*discordData
}

func (b ArkwBot) Implementation() string {
	return "Arikawa"
}

func (b ArkwBot) Init(token string) error {
	b.token = token
	b.s = state.New("Bot " + b.token)
	return nil
}

func (b ArkwBot) connect() error {
	now := time.Now()
	err := b.s.Open(context.Background())
	if err != nil {
		return err
	}
	b.startedAt = &now
	return nil
}

func (b ArkwBot) Start() error {
	b.registerListeners()
	return b.connect()
}

func (b ArkwBot) StartedAt() *time.Time {
	return b.startedAt
}

func (b ArkwBot) CountUsersInVoiceChannel(ch model.VoiceChannel) (count int) {
	sf, err := dc.ParseSnowflake(ch.Guild().ID())
	if err != nil {
		return
	}
	states, err := b.s.VoiceStates(dc.GuildID(sf))
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

func (b ArkwBot) Stop() error {
	return b.s.Close()
}

func (b ArkwBot) Self() (model.User, error) {
	u, err := b.s.Me()
	if err != nil {
		return nil, err
	}
	return buildUser(u.ID.String()), nil
}

func (b ArkwBot) GuildCount() int {
	v, _ := b.s.Guilds()
	return len(v)
}

func (b ArkwBot) RegisterSlashCommands() error {
	return registerCommands(b)
}

func (b ArkwBot) SendComplexMessage(channelID string, message *model.ComplexMessage) (model.Message, error) {
	var refMessage *dc.MessageReference

	channelSf, err := dc.ParseSnowflake(channelID)
	if err != nil {
		return nil, err
	}

	if message.ReplyTo != nil {
		refMessageSf, err := dc.ParseSnowflake(message.ReplyTo.ID())
		if err != nil {
			return nil, err
		}

		var guildSf dc.Snowflake

		if message.ReplyTo.Channel().Type() == model.ChannelTypeGuild {
			guildSf, err = dc.ParseSnowflake(message.ReplyTo.Channel().Guild().ID())
			if err != nil {
				return nil, err
			}
		}

		refMessage = &dc.MessageReference{
			MessageID: dc.MessageID(refMessageSf),
			ChannelID: dc.ChannelID(channelSf),
			GuildID:   dc.GuildID(guildSf),
		}
	}

	arkwMessage := prepareComplexMessage(message)

	msg, err := b.s.SendMessageComplex(dc.ChannelID(channelSf), api.SendMessageData{
		Content:    message.Content,
		Embeds:     arkwMessage.embed,
		Reference:  refMessage,
		Components: arkwMessage.components,
	})
	if err != nil {
		return nil, err
	}

	cType := model.ChannelTypeGuild
	if msg.GuildID.String() == "" {
		cType = model.ChannelTypeDirect
	}

	return buildMessage(
		msg.ID.String(), buildChannel(channelID, buildGuild(msg.GuildID.String()), cType),
		buildUser(msg.Author.ID.String()),
		message.Content,
	), nil
}

func (b ArkwBot) EditComplexMessage(m model.Message, newMessage *model.ComplexMessage) (model.Message, error) {
	oldMsgSf, err := dc.ParseSnowflake(m.ID())
	if err != nil {
		return nil, err
	}

	channelSf, err := dc.ParseSnowflake(m.Channel().ID())
	if err != nil {
		return nil, err
	}

	arkwMessage := prepareComplexMessage(newMessage)

	var embedsPtr *[]dc.Embed
	var componentsPtr *dc.ContainerComponents

	if len(arkwMessage.embed) > 0 {
		embedsPtr = &arkwMessage.embed
	}

	if len(arkwMessage.components) > 0 {
		componentsPtr = &arkwMessage.components
	}

	msg, err := b.s.EditMessageComplex(dc.ChannelID(channelSf), dc.MessageID(oldMsgSf), api.EditMessageData{
		Content:    option.NewNullableString(newMessage.Content),
		Embeds:     embedsPtr,
		Components: componentsPtr,
	})
	if err != nil {
		return nil, err
	}

	return buildMessage(
		msg.ID.String(), buildChannel(m.Channel().ID(), buildGuild(msg.GuildID.String()), m.Channel().Type()),
		buildUser(msg.Author.ID.String()),
		newMessage.Content,
	), nil
}

func (b ArkwBot) SendReplyMessage(m model.Message, content string) (model.Message, error) {
	return b.SendComplexMessage(m.Channel().ID(), &model.ComplexMessage{
		Content: content,
		ReplyTo: m,
	})
}

func (b ArkwBot) SendMessage(channelID string, message string) (model.Message, error) {
	return b.SendComplexMessage(channelID, &model.ComplexMessage{
		Content: message,
	})
}

func (b ArkwBot) SendReplyEmbedMessage(m model.Message, embed *model.Embed) (model.Message, error) {
	return b.SendComplexMessage(m.Channel().ID(), &model.ComplexMessage{
		Embeds:  []*model.Embed{embed},
		ReplyTo: m,
	})
}

func (b ArkwBot) SendEmbedMessage(channelID string, embed *model.Embed) (model.Message, error) {
	return b.SendComplexMessage(channelID, &model.ComplexMessage{
		Embeds: []*model.Embed{embed},
	})
}

func (b ArkwBot) EditMessageContent(message model.Message, newContent string) (model.Message, error) {
	return b.EditComplexMessage(message, &model.ComplexMessage{
		Content: newContent,
	})
}

func (b ArkwBot) EditMessageEmbed(message model.Message, newEmbed *model.Embed) (model.Message, error) {
	return b.EditComplexMessage(message, &model.ComplexMessage{
		Embeds: []*model.Embed{newEmbed},
	})
}

func (b ArkwBot) OpenChannelWithUser(userID string) (model.TextChannel, error) {
	sf, err := dc.ParseSnowflake(userID)
	if err != nil {
		return nil, err
	}
	dm, err := b.s.CreatePrivateChannel(dc.UserID(sf))
	if err != nil {
		return nil, err
	}
	return buildChannel(dm.ID.String(), buildGuild(""), model.ChannelTypeDirect), nil
}

func (b ArkwBot) IsLive() bool {
	_, err := b.s.Me()
	return err == nil
}

func (b ArkwBot) Latency() time.Duration {
	return b.s.Gateway().Latency()
}

func (b ArkwBot) OpenGuild(guildID string) (model.Guild, error) {
	sf, err := dc.ParseSnowflake(guildID)
	if err != nil {
		return nil, err
	}
	g, err := b.s.Guild(dc.GuildID(sf))
	if err != nil {
		return nil, err
	}
	return buildGuild(g.ID.String()), nil
}

func (b ArkwBot) JoinVoiceChannel(guildID, channelID string) (model.VoiceConnection, error) {
	vs, err := voice.NewSession(b.s)
	if err != nil {
		return nil, err
	}
	channelSf, err := dc.ParseSnowflake(channelID)
	if err != nil {
		return nil, err
	}
	err = vs.JoinChannel(context.Background(), dc.ChannelID(channelSf), false, true)
	if err != nil {
		return nil, err
	}
	return buildVoiceConnection(vs), nil
}

func (b ArkwBot) FindUserVoiceState(guildID, userID string) (model.VoiceState, error) {
	guildSf, err := dc.ParseSnowflake(guildID)
	if err != nil {
		return nil, err
	}
	userSf, err := dc.ParseSnowflake(userID)
	if err != nil {
		return nil, err
	}
	vs, err := b.s.VoiceState(dc.GuildID(guildSf), dc.UserID(userSf))
	if err != nil {
		return nil, err
	}
	return buildVoiceState(buildVoiceChannel(vs.ChannelID.String(), buildGuild(guildID))), nil
}

func (b ArkwBot) GetMember(guildID, channelID, userID string) (model.Member, error) {
	if guildID == "" {
		return nil, errors.New("guildID is empty")
	}

	guildSf, err := dc.ParseSnowflake(guildID)
	if err != nil {
		return nil, err
	}

	channelSf, err := dc.ParseSnowflake(channelID)
	if err != nil && channelID != "" {
		return nil, err
	}

	userSf, err := dc.ParseSnowflake(userID)
	if err != nil {
		return nil, err
	}

	guild, err := b.s.Guild(dc.GuildID(guildSf))
	if err != nil {
		return nil, err
	}

	var channel *dc.Channel

	if channelID == "" {
		channel = new(dc.Channel)
	} else {
		channel, err = b.s.Channel(dc.ChannelID(channelSf))
		if err != nil {
			return nil, err
		}
	}

	m, err := b.s.Member(dc.GuildID(guildSf), dc.UserID(userSf))
	if err != nil {
		return nil, err
	}

	perms := dc.CalcOverwrites(
		*guild,
		*channel,
		*m,
	)

	var roles []model.Role
	for _, r := range m.RoleIDs {
		role, err := b.s.Role(dc.GuildID(guildSf), r)
		if err != nil {
			return nil, err
		}
		roles = append(roles, buildRole(role))
	}

	return buildMember(
		roles,
		model.Permissions(uint64(perms)),
	), nil
}

func (b ArkwBot) UpdatePresence(presence *model.Presence) error {
	var ty dc.ActivityType
	switch presence.Type {
	case model.PresencePlaying:
		ty = dc.GameActivity
	case model.PresenceListening:
		ty = dc.ListeningActivity
	case model.PresenceStreaming:
		ty = dc.StreamingActivity
	default:
		return errors.New("invalid presence type")
	}
	return b.s.Gateway().Send(context.Background(), &gateway.UpdatePresenceCommand{
		Activities: []dc.Activity{
			{Name: presence.Title, URL: presence.Extra, Type: ty},
		},
	})
}
