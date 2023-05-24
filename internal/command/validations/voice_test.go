package validations_test

import (
	"os"
	"testing"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/stretchr/testify/require"
)

var (
	defaultLang *i18n.Language
)

func TestMain(m *testing.M) {
	var err error
	i18n.I18nRootDir = "../../../assets/i18n"
	defaultLang, err = i18n.GetLanguage(i18n.DefaultLanguageName)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestMustHaveVoicerOnGuild(t *testing.T) {
	t.Run("implementation is missing", func(t *testing.T) {
		ctx := &command.CommandContext{
			Lang: defaultLang,
		}
		b, _ := validations.MustHaveVoicerOnGuild.Checker(ctx)
		require.False(t, b)
	})

	t.Run("hasn't voicer in the same guild id", func(t *testing.T) {
		discord.UseImplementation(discord.DummyBot{})

		guildId := "321"

		v, err := voicer.NewVoicerForUser("123", guildId)
		require.NotNil(t, v)
		require.Nil(t, err)

		ctx := &command.CommandContext{
			GuildID: "123123", // not the same guild id
			Locals:  make(map[string]interface{}),
			Lang:    defaultLang,
		}

		b, _ := validations.MustHaveVoicerOnGuild.Checker(ctx)
		require.False(t, b)
	})

	t.Run("has voicer in the same guild id", func(t *testing.T) {
		discord.UseImplementation(discord.DummyBot{})

		guildId := "321"

		v, err := voicer.NewVoicerForUser("123", guildId)
		require.NotNil(t, v)
		require.Nil(t, err)

		ctx := &command.CommandContext{
			GuildID: guildId,
			Locals:  make(map[string]interface{}),
			Lang:    defaultLang,
		}

		b, _ := validations.MustHaveVoicerOnGuild.Checker(ctx)
		require.True(t, b)
	})
}
