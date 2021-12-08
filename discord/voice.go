package discord

type VoiceState interface {
	Channel() VoiceChannel
}

type VoiceConnection interface {
	OpusSend() chan []byte
	Speaking(bool) error
	Disconnect() error
}
