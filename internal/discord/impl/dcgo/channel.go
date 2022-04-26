package dcgo

import "github.com/Pauloo27/aryzona/internal/discord/model"

type Channel struct {
	id    string
	guild Guild
	cType model.ChannelType
}

func (c Channel) ID() string {
	return c.id
}

func (c Channel) Guild() model.Guild {
	return c.guild
}

func (c Channel) Type() model.ChannelType {
	return c.cType
}

func buildChannel(id string, guild Guild, cType model.ChannelType) Channel {
	return Channel{
		id:    id,
		guild: guild,
		cType: cType,
	}
}
