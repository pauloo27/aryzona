package discord

type VoiceState interface {
	Channel() VoiceChannel
}

type VoiceConnection interface {
	WriteOpus([]byte) error
	Speaking(bool) error
	Disconnect() error
}
