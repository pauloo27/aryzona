package discord

type Channel interface {
	ID() string
	Guild() Guild
}

type VoiceChannel interface {
	Channel
}
