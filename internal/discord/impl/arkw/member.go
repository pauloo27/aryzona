package arkw

import "github.com/Pauloo27/aryzona/internal/discord/model"

type Member struct {
	roles []model.Role
	perms model.Permissions
}

func (m Member) Roles() []model.Role {
	return m.roles
}

func (m Member) Permissions() model.Permissions {
	return m.perms
}

func buildMember(roles []model.Role, perms model.Permissions) Member {
	return Member{
		roles: roles,
		perms: perms,
	}
}
