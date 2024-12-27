package disgo

type Guild struct {
	id string
}

func (g Guild) ID() string {
	return g.id
}

func buildGuild(id string) Guild {
	return Guild{
		id: id,
	}
}
