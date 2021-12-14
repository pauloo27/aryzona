package dcgo

import "github.com/Pauloo27/aryzona/internal/discord"

type Channel struct {
	id    string
	guild Guild
}

func (c Channel) ID() string {
	return c.id
}

func (c Channel) Guild() discord.Guild {
	return c.guild
}

func buildChannel(id string, guild Guild) Channel {
	return Channel{
		id:    id,
		guild: guild,
	}
}
