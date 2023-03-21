package discord

import (
	"time"

	"github.com/Pauloo27/aryzona/internal/discord/event"
	"github.com/Pauloo27/aryzona/internal/discord/model"
)

var (
	Bot BotAdapter
)

type BotAdapter interface {
	Implementation() string
	Init(token string) error
	StartedAt() *time.Time
	Listen(event event.EventType, handlerFunc any) error
	Start() error
	Stop() error
	Self() (model.User, error)
	GetMember(guildID, channelID, userID string) (model.Member, error)
	CountUsersInVoiceChannel(vc model.VoiceChannel) int
	SendMessage(channelID string, content string) (model.Message, error)
	SendComplexMessage(channelID string, message *model.ComplexMessage) (model.Message, error)
	EditComplexMessage(message model.Message, newMessage *model.ComplexMessage) (model.Message, error)
	SendReplyMessage(message model.Message, content string) (model.Message, error)
	SendReplyEmbedMessage(message model.Message, embed *model.Embed) (model.Message, error)
	SendEmbedMessage(channelID string, embed *model.Embed) (model.Message, error)
	EditMessageContent(message model.Message, newContent string) (model.Message, error)
	EditMessageEmbed(message model.Message, embed *model.Embed) (model.Message, error)
	OpenChannelWithUser(userID string) (model.TextChannel, error)
	OpenGuild(guildID string) (model.Guild, error)
	Latency() time.Duration
	JoinVoiceChannel(guildID, channelID string) (model.VoiceConnection, error)
	FindUserVoiceState(guildID string, userID string) (model.VoiceState, error)
	UpdatePresence(presence *model.Presence) error
	GuildCount() int
	RegisterSlashCommands() error
}

func UseImplementation(bot BotAdapter) {
	Bot = bot
}

func CreateBot(token string) error {
	return Bot.Init(token)
}
