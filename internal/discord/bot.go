package discord

import (
	"time"

	"github.com/Pauloo27/aryzona/internal/discord/event"
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
	Self() (User, error)
	CountUsersInVoiceChannel(vc VoiceChannel) int
	SendMessage(channelID string, content string) (Message, error)
	SendReplyMessage(message Message, content string) (Message, error)
	SendReplyEmbedMessage(message Message, embed *Embed) (Message, error)
	SendEmbedMessage(channelID string, embed *Embed) (Message, error)
	OpenChannelWithUser(userID string) (Channel, error)
	OpenGuild(guildID string) (Guild, error)
	Latency() time.Duration
	JoinVoiceChannel(guildID, channelID string) (VoiceConnection, error)
	FindUserVoiceState(guildID string, userID string) (VoiceState, error)
	UpdatePresence(presence *Presence) error
	GuildCount() int
	RegisterSlashCommands() error
}

func UseImplementation(bot BotAdapter) {
	Bot = bot
}

func CreateBot(token string) error {
	return Bot.Init(token)
}
