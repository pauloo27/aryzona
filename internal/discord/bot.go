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
	Listen(event event.EventType, handlerFunc interface{}) error
	Start() error
	Stop() error
	Self() (model.User, error)
	CountUsersInVoiceChannel(vc model.VoiceChannel) int
	SendMessage(channelID string, content string) (model.Message, error)
	SendReplyMessage(message model.Message, content string) (model.Message, error)
	SendReplyEmbedMessage(message model.Message, embed *Embed) (model.Message, error)
	SendEmbedMessage(channelID string, embed *Embed) (model.Message, error)
	OpenChannelWithUser(userID string) (model.Channel, error)
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
