package model

type ChannelType string

const (
	ChannelTypeGuild  ChannelType = "GUILD"
	ChannelTypeDirect ChannelType = "DIRECT"
)

type TextChannel interface {
	ID() string
	Guild() Guild
	Type() ChannelType
}

type VoiceChannel interface {
	ID() string
	Guild() Guild
}
