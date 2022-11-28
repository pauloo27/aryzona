package arkw

import (
	"context"

	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
)

type VoiceChannel struct {
	id    string
	guild Guild
}

func (c VoiceChannel) ID() string {
	return c.id
}

func (c VoiceChannel) Guild() model.Guild {
	return c.guild
}

func buildVoiceChannel(id string, guild Guild) VoiceChannel {
	return VoiceChannel{
		id:    id,
		guild: guild,
	}
}

type VoiceState struct {
	channel VoiceChannel
}

func (c VoiceState) Channel() model.VoiceChannel {
	return c.channel
}

func buildVoiceState(channel VoiceChannel) VoiceState {
	return VoiceState{
		channel: channel,
	}
}

type VoiceConnection struct {
	session *voice.Session
}

func (c VoiceConnection) WriteOpus(b []byte) (int, error) {
	return c.session.Write(b)
}

func (c VoiceConnection) Speaking(speaking bool) error {
	if speaking {
		return c.session.Speaking(context.Background(), voicegateway.Microphone)
	}
	return c.session.Speaking(context.Background(), voicegateway.NotSpeaking)
}

func (c VoiceConnection) Disconnect() error {
	return c.session.Leave(context.Background())
}

func buildVoiceConnection(session *voice.Session) VoiceConnection {
	return VoiceConnection{
		session: session,
	}
}
