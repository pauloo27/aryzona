package dcgo

import "github.com/Pauloo27/aryzona/internal/discord/model"

type Member struct {
	roles []model.Role
}

func (m Member) Roles() []model.Role {
	return m.roles
}

func buildMember(roles []model.Role) Member {
	return Member{
		roles: roles,
	}
}
