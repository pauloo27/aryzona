package model

type Channel interface {
	ID() string
	Guild() Guild
}

type VoiceChannel interface {
	Channel
}
