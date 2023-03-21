package permissions_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/permissions"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"

	"github.com/stretchr/testify/require"
)

func TestAdminPermission(t *testing.T) {
	member := &DummyMember{
		permissions: model.PermissionSendMessages,
	}

	adminMember := &DummyMember{
		permissions: model.PermissionAdministrator,
	}

	memberCtx := &command.CommandContext{
		GuildID: "123",
		Bot: &DummyBot{member: member},
	}

	adminMemberCtx := &command.CommandContext{
		GuildID: "123",
		Bot: &DummyBot{member: adminMember},
	}

	dmCtx := &command.CommandContext{
		GuildID: "",
	}

	require.False(t, permissions.MustBeAdmin.Checker(memberCtx))

	require.True(t, permissions.MustBeAdmin.Checker(adminMemberCtx))

	require.False(t, permissions.MustBeAdmin.Checker(dmCtx))
}

type DummyMember struct {
	roles       []model.Role
	permissions model.Permissions
}

func (m DummyMember) Roles() []model.Role {
	return m.roles
}

func (m DummyMember) Permissions() model.Permissions {
	return m.permissions
}

type DummyBot struct {
	discord.DummyBot

	member *DummyMember
}

func (b DummyBot) GetMember(guildId, userId, authorId string) (model.Member, error) {
	return b.member, nil
}
