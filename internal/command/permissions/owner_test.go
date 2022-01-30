package permissions_test

import (
	"os"
	"testing"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/permissions"

	"github.com/stretchr/testify/require"
)

func TestOwnerPermission(t *testing.T) {
	ctx := &command.CommandContext{
		AuthorID: "777",
	}

	os.Setenv("DC_BOT_OWNER_ID", "123")
	require.False(t, permissions.MustBeOwner.Checker(ctx))

	os.Setenv("DC_BOT_OWNER_ID", "")
	require.False(t, permissions.MustBeOwner.Checker(ctx))

	os.Setenv("DC_BOT_OWNER_ID", "321")
	require.False(t, permissions.MustBeOwner.Checker(ctx))

	os.Setenv("DC_BOT_OWNER_ID", "7 7 7")
	require.False(t, permissions.MustBeOwner.Checker(ctx))

	os.Setenv("DC_BOT_OWNER_ID", "777")
	require.True(t, permissions.MustBeOwner.Checker(ctx))
}
