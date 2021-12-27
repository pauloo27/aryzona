package dcgo

import (
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/bwmarrin/discordgo"
)

type VoiceConnection struct {
	connection *discordgo.VoiceConnection
}

func (c VoiceConnection) WriteOpus(b []byte) (int, error) {
	c.connection.OpusSend <- b
	return len(b), nil
}

func (c VoiceConnection) Speaking(speaking bool) error {
	return c.connection.Speaking(speaking)
}

func (c VoiceConnection) Disconnect() error {
	return c.connection.Disconnect()
}

func buildVoiceConnection(connection *discordgo.VoiceConnection) VoiceConnection {
	return VoiceConnection{
		connection: connection,
	}
}

type VoiceChannel struct {
	id    string
	guild Guild
}

func (c VoiceChannel) ID() string {
	return c.id
}

func (c VoiceChannel) Guild() discord.Guild {
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

func (c VoiceState) Channel() discord.VoiceChannel {
	return c.channel
}

func buildVoiceState(channel VoiceChannel) VoiceState {
	return VoiceState{
		channel: channel,
	}
}
