package discord

type VoiceState interface {
	ChanID() string
	Speaking(bool) error
	Disconnect() error
	Connection() VoiceConnection
}

type VoiceConnection interface {
	OpusSend() chan []byte
}
