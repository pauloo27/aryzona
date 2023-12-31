package permissions_test

import (
	"testing"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/permissions"
	"github.com/pauloo27/aryzona/internal/config"

	"github.com/stretchr/testify/require"
)

func TestOwnerPermission(t *testing.T) {
	ctx := &command.Context{
		AuthorID: "777",
	}

	config.Config.OwnerID = "123"
	require.False(t, permissions.MustBeOwner.Checker(ctx))

	config.Config.OwnerID = ""
	require.False(t, permissions.MustBeOwner.Checker(ctx))

	config.Config.OwnerID = "321"
	require.False(t, permissions.MustBeOwner.Checker(ctx))

	config.Config.OwnerID = "7 7 7"
	require.False(t, permissions.MustBeOwner.Checker(ctx))

	config.Config.OwnerID = "777"
	require.True(t, permissions.MustBeOwner.Checker(ctx))
}
