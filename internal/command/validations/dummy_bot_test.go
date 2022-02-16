package validations_test

import (
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/event"
)

type DummyBot struct {
}

var _ discord.BotAdapter = DummyBot{}

func (DummyBot) Implementation() string {
	return "Dummy Bot"
}

func (DummyBot) Init(token string) error {
	return nil
}

func (DummyBot) StartedAt() *time.Time {
	return nil
}

func (DummyBot) Listen(event event.EventType, handlerFunc interface{}) error {
	return nil
}

func (DummyBot) Start() error {
	return nil
}

func (DummyBot) Stop() error {
	return nil
}

func (DummyBot) Self() (discord.User, error) {
	return nil, nil
}

func (DummyBot) CountUsersInVoiceChannel(vc discord.VoiceChannel) int {
	return 0
}

func (DummyBot) SendMessage(channelID string, content string) (discord.Message, error) {
	return nil, nil
}

func (DummyBot) SendReplyMessage(message discord.Message, content string) (discord.Message, error) {
	return nil, nil
}

func (DummyBot) SendReplyEmbedMessage(message discord.Message, embed *discord.Embed) (discord.Message, error) {
	return nil, nil
}

func (DummyBot) SendEmbedMessage(channelID string, embed *discord.Embed) (discord.Message, error) {
	return nil, nil
}

func (DummyBot) OpenChannelWithUser(userID string) (discord.Channel, error) {
	return nil, nil
}
func (DummyBot) OpenGuild(guildID string) (discord.Guild, error) {
	return Guild{}, nil
}

func (DummyBot) Latency() time.Duration {
	return time.Second * 20
}

func (DummyBot) JoinVoiceChannel(guildID, channelID string) (discord.VoiceConnection, error) {
	return nil, nil
}

func (DummyBot) FindUserVoiceState(guildID string, userID string) (discord.VoiceState, error) {
	return VoiceState{}, nil
}

func (DummyBot) UpdatePresence(presence *discord.Presence) error {
	return nil
}

func (DummyBot) GuildCount() int {
	return 0
}

func (DummyBot) RegisterSlashCommands() error {
	return nil
}

type VoiceState struct {
}

func (VoiceState) Channel() discord.VoiceChannel {
	return VoiceChannel{}
}

type Guild struct {
}

func (Guild) ID() string {
	return "1233"
}

type VoiceChannel struct {
}

func (VoiceChannel) ID() string {
	return "12335"
}

func (VoiceChannel) Guild() discord.Guild {
	return Guild{}
}
