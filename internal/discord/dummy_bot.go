package discord

import (
	"time"

	"github.com/pauloo27/aryzona/internal/discord/event"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

type DummyBot struct {
}

var _ BotAdapter = DummyBot{}

func (DummyBot) EditMessageContent(message model.Message, newContent string) (model.Message, error) {
	return nil, nil
}

func (DummyBot) EditComplexMessage(message model.Message, newMessage *model.ComplexMessage) (model.Message, error) {
	return nil, nil
}

func (DummyBot) EditMessageEmbed(message model.Message, embed *model.Embed) (model.Message, error) {
	return nil, nil
}

func (DummyBot) Implementation() string {
	return "Dummy Bot"
}

func (DummyBot) GetMember(guildID, channelID, memberID string) (model.Member, error) {
	return nil, nil
}

func (DummyBot) Init(token string) error {
	return nil
}

func (DummyBot) StartedAt() *time.Time {
	return nil
}

func (DummyBot) Listen(event event.EventType, handlerFunc any) error {
	return nil
}

func (DummyBot) Start() error {
	return nil
}

func (DummyBot) Stop() error {
	return nil
}

func (DummyBot) Self() (model.User, error) {
	return nil, nil
}

func (DummyBot) CountUsersInVoiceChannel(vc model.VoiceChannel) int {
	return 0
}

func (DummyBot) SendMessage(channelID string, content string) (model.Message, error) {
	return nil, nil
}

func (DummyBot) SendReplyMessage(message model.Message, content string) (model.Message, error) {
	return nil, nil
}

func (DummyBot) SendReplyEmbedMessage(message model.Message, embed *model.Embed) (model.Message, error) {
	return nil, nil
}

func (DummyBot) SendEmbedMessage(channelID string, embed *model.Embed) (model.Message, error) {
	return nil, nil
}

func (DummyBot) SendComplexMessage(channelID string, message *model.ComplexMessage) (model.Message, error) {
	return nil, nil
}

func (DummyBot) OpenChannelWithUser(userID string) (model.TextChannel, error) {
	return nil, nil
}
func (DummyBot) OpenGuild(guildID string) (model.Guild, error) {
	return Guild{}, nil
}

func (DummyBot) Latency() time.Duration {
	return time.Second * 20
}

func (DummyBot) JoinVoiceChannel(guildID, channelID string) (model.VoiceConnection, error) {
	return nil, nil
}

func (DummyBot) FindUserVoiceState(guildID string, userID string) (model.VoiceState, error) {
	return VoiceState{}, nil
}

func (DummyBot) UpdatePresence(presence *model.Presence) error {
	return nil
}

func (DummyBot) IsLive() bool {
	return true
}

func (DummyBot) GuildCount() int {
	return 0
}

func (DummyBot) RegisterSlashCommands() error {
	return nil
}

type VoiceState struct {
}

func (VoiceState) Channel() model.VoiceChannel {
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

func (VoiceChannel) Guild() model.Guild {
	return Guild{}
}
