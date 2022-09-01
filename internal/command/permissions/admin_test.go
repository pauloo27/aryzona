package permissions_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/permissions"
	"github.com/Pauloo27/aryzona/internal/discord/model"

	"github.com/stretchr/testify/require"
)

func TestAdminPermission(t *testing.T) {
	adminCtx := &command.CommandContext{
		GuildID: "123123",
		Member: &DummyMember{
			roles: []model.Role{
				DummyRole{
					permissions: model.PermissionManageMessages | model.PermissionAdministrator,
				},
			},
		},
	}

	multRolesAdmin := &command.CommandContext{
		GuildID: "123123",
		Member: &DummyMember{
			roles: []model.Role{
				DummyRole{
					permissions: model.PermissionManageMessages,
				},
				DummyRole{
					permissions: model.PermissionAdministrator,
				},
			},
		},
	}

	memberCtx := &command.CommandContext{
		GuildID: "123123",
		Member: &DummyMember{
			roles: []model.Role{
				DummyRole{
					permissions: model.PermissionSendMessages,
				},
			},
		},
	}

	dmCtx := &command.CommandContext{}

	require.True(t, permissions.MustBeAdmin.Checker(adminCtx))
	require.True(t, permissions.MustBeAdmin.Checker(multRolesAdmin))
	require.False(t, permissions.MustBeAdmin.Checker(memberCtx))
	require.False(t, permissions.MustBeAdmin.Checker(dmCtx))
}

type DummyRole struct {
	permissions model.Permissions
}

type DummyMember struct {
	roles []model.Role
}

func (m DummyMember) Roles() []model.Role {
	return m.roles
}

func (r DummyRole) Color() int {
	return 0
}

func (r DummyRole) Permissions() model.Permissions {
	return r.permissions
}

func (r DummyRole) ID() string {
	return "dummy"
}

func (r DummyRole) Mentionable() bool {
	return false
}

func (r DummyRole) Name() string {
	return "Dummy"
}

func (r DummyRole) Position() int {
	return 0
}
