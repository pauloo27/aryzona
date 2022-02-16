package validations_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/stretchr/testify/require"
)

func TestMustHaveVoicerOnGuild(t *testing.T) {
	t.Run("implementation is missing", func(t *testing.T) {
		ctx := &command.CommandContext{}
		b, _ := validations.MustHaveVoicerOnGuild.Checker(ctx)
		require.False(t, b)
	})

	t.Run("hasn't voicer in the same guild id", func(t *testing.T) {
		discord.UseImplementation(DummyBot{})

		guildId := "321"

		v, err := voicer.NewVoicerForUser("123", guildId)
		require.NotNil(t, v)
		require.Nil(t, err)

		ctx := &command.CommandContext{
			GuildID: "123123", // not the same guild id
			Locals:  make(map[string]interface{}),
		}

		b, _ := validations.MustHaveVoicerOnGuild.Checker(ctx)
		require.False(t, b)
	})

	t.Run("has voicer in the same guild id", func(t *testing.T) {
		discord.UseImplementation(DummyBot{})

		guildId := "321"

		v, err := voicer.NewVoicerForUser("123", guildId)
		require.NotNil(t, v)
		require.Nil(t, err)

		ctx := &command.CommandContext{
			GuildID: guildId,
			Locals:  make(map[string]interface{}),
		}

		b, _ := validations.MustHaveVoicerOnGuild.Checker(ctx)
		require.True(t, b)
	})
}
