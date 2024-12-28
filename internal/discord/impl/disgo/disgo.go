package disgo

import (
	"context"
	"errors"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	dc "github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

var _ discord.BotAdapter = &DisgoBot{}

type DisgoBot struct {
	token     string
	startedAt time.Time
	client    bot.Client
	opts      []bot.ConfigOpt
	intents   []gateway.Intents
}

func (d *DisgoBot) CountUsersInVoiceChannel(vc model.VoiceChannel) int {
	panic("unimplemented")
}

func (d *DisgoBot) EditComplexMessage(message model.Message, newMessage *model.ComplexMessage) (model.Message, error) {
	panic("unimplemented")
}

func (d *DisgoBot) EditMessageContent(message model.Message, newContent string) (model.Message, error) {
	panic("unimplemented")
}

func (d *DisgoBot) EditMessageEmbed(message model.Message, embed *model.Embed) (model.Message, error) {
	panic("unimplemented")
}

func (d *DisgoBot) FindUserVoiceState(guildID string, userID string) (model.VoiceState, error) {
	panic("unimplemented")
}

func (d *DisgoBot) GetMember(guildID string, channelID string, userID string) (model.Member, error) {
	panic("unimplemented")
}

func (d *DisgoBot) GuildCount() int {
	// TODO: feels stupid 
	guilds, err := d.client.Rest().GetCurrentUserGuilds("", snowflake.MustParse("null"), snowflake.MustParse("null"), 100, false)
	if err != nil {
		return -1
	}
	return len(guilds)
}

func (d *DisgoBot) Implementation() string {
	return "Disgo"
}

func (d *DisgoBot) Init(token string) error {
	d.token = token
	d.intents = append(d.intents, gateway.IntentGuilds)
	return nil
}

func (d *DisgoBot) IsLive() bool {
	_, err := d.Self()
	return err == nil
}

func (d *DisgoBot) JoinVoiceChannel(guildID string, channelID string) (model.VoiceConnection, error) {
	panic("unimplemented")
}

func (d *DisgoBot) Latency() time.Duration {
	panic("unimplemented")
}

func (d *DisgoBot) OpenChannelWithUser(userID string) (model.TextChannel, error) {
	userSf, err := snowflake.Parse(userID)
	if err != nil {
		return nil, err
	}
	dm, err := d.client.Rest().CreateDMChannel(userSf)
	if err != nil {
		return nil, err
	}

	return buildChannel(dm.ID().String(), buildGuild(""), model.ChannelTypeDirect), nil
}

func (d *DisgoBot) OpenGuild(guildID string) (model.Guild, error) {
	panic("unimplemented")
}

func (d *DisgoBot) RegisterSlashCommands() error {
	// TODO:
	return nil
}

func (d *DisgoBot) Self() (model.User, error) {
	return buildUser(d.client.ID().String()), nil
}

func (b *DisgoBot) SendReplyMessage(m model.Message, content string) (model.Message, error) {
	return b.SendComplexMessage(m.Channel().ID(), &model.ComplexMessage{
		Content: content,
		ReplyTo: m,
	})
}

func (b *DisgoBot) SendMessage(channelID string, message string) (model.Message, error) {
	return b.SendComplexMessage(channelID, &model.ComplexMessage{
		Content: message,
	})
}

func (b *DisgoBot) SendReplyEmbedMessage(m model.Message, embed *model.Embed) (model.Message, error) {
	return b.SendComplexMessage(m.Channel().ID(), &model.ComplexMessage{
		Embeds:  []*model.Embed{embed},
		ReplyTo: m,
	})
}

func (b *DisgoBot) SendEmbedMessage(channelID string, embed *model.Embed) (model.Message, error) {
	return b.SendComplexMessage(channelID, &model.ComplexMessage{
		Embeds: []*model.Embed{embed},
	})
}

func (d *DisgoBot) SendComplexMessage(channelID string, message *model.ComplexMessage) (model.Message, error) {
	channelSf, err := snowflake.Parse(channelID)
	if err != nil {
		return nil, err
	}

	messageBuilder := dc.NewMessageCreateBuilder()

	if message.ReplyTo != nil {
		replyToSf, err := snowflake.Parse(message.ReplyTo.ID())
		if err != nil {
			return nil, err
		}
		messageBuilder.SetMessageReferenceByID(replyToSf)
	}

	if message.Content != "" {
		messageBuilder.SetContent(message.Content)
	}

	if len(message.Embeds) > 0 {
		for _, embed := range message.Embeds {
			messageBuilder.AddEmbeds(buildEmbed(embed))
		}
	}

	if len(message.ComponentRows) > 0 {
		// TODO:
	}

	sentMessage, err := d.client.Rest().CreateMessage(channelSf, messageBuilder.Build())
	if err != nil {
		return nil, err
	}

	channelType := model.ChannelTypeDirect
	guildID := ""

	if sentMessage.GuildID != nil {
		channelType = model.ChannelTypeGuild
		guildID = sentMessage.GuildID.String()
	}

	return buildMessage(
		sentMessage.ID.String(),
		buildChannel(sentMessage.ChannelID.String(), buildGuild(guildID), channelType),
		buildUser(sentMessage.Author.ID.String()),
		sentMessage.Content,
	), nil
}

func (d *DisgoBot) Start() error {
	opts := append(d.opts, bot.WithGatewayConfigOpts(gateway.WithIntents(d.intents...)))

	client, err := disgo.New(
		d.token,
		opts...,
	)
	if err != nil {
		return err
	}
	d.startedAt = time.Now()
	d.client = client

	return d.client.OpenGateway(context.Background())
}

func (d *DisgoBot) StartedAt() *time.Time {
	if d.startedAt.IsZero() {
		return nil
	}
	return &d.startedAt
}

func (d *DisgoBot) Stop() error {
	d.client.Close(context.Background())
	return nil
}

func (d *DisgoBot) UpdatePresence(presence *model.Presence) error {
	var opt gateway.PresenceOpt
	switch presence.Type {
	case model.PresencePlaying:
		opt = gateway.WithPlayingActivity(presence.Title)
	case model.PresenceListening:
		opt = gateway.WithListeningActivity(presence.Title)
	case model.PresenceStreaming:
		opt = gateway.WithStreamingActivity(presence.Title, presence.Extra)
	default:
		return errors.New("invalid presence type")
	}
	return d.client.SetPresence(context.Background(), opt)
}
