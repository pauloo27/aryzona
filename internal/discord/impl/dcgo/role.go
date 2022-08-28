package dcgo

import (
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/bwmarrin/discordgo"
)

type Role struct {
	id, name        string
	permissions     model.Permissions
	position, color int
	mentionable     bool
}

func (r Role) ID() string {
	return r.id
}

func (r Role) Name() string {
	return r.name
}

func (r Role) Permissions() model.Permissions {
	return r.permissions
}

func (r Role) Position() int {
	return r.position
}

func (r Role) Color() int {
	return r.color
}

func (r Role) Mentionable() bool {
	return r.mentionable
}

func buildRole(role *discordgo.Role) Role {
	return Role{
		id:          role.ID,
		name:        role.Name,
		permissions: model.Permissions(role.Permissions),
		position:    role.Position,
		color:       int(role.Color),
		mentionable: role.Mentionable,
	}
}
