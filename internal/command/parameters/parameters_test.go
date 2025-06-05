package parameters_test

import (
	"os"
	"testing"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/stretchr/testify/require"
)

var (
	defaultLang *i18n.Language
)

func TestMain(m *testing.M) {
	i18n.I18nRootDir = "../../../assets/i18n"
	loadDefaultLang()
	os.Exit(m.Run())
}

func TestBoolean(t *testing.T) {
	var b any
	var err error

	ctx := &command.Context{
		Lang: defaultLang,
	}

	b, err = parameters.ParameterBool.Parser(ctx, 0, []string{"true"})
	require.Nil(t, err)
	require.Equal(t, b, true)

	b, err = parameters.ParameterBool.Parser(ctx, 1, []string{"hello", "true"})
	require.Nil(t, err)
	require.Equal(t, b, true)

	b, err = parameters.ParameterBool.Parser(ctx, 0, []string{"false"})
	require.Nil(t, err)
	require.Equal(t, b, false)

	b, err = parameters.ParameterBool.Parser(ctx, 1, []string{"hello", "false"})
	require.Nil(t, err)
	require.Equal(t, b, false)

	_, err = parameters.ParameterBool.Parser(ctx, 0, []string{"invalid"})
	require.NotNil(t, err)
}

func TestInt(t *testing.T) {
	var i any
	var err error

	ctx := &command.Context{
		Lang: defaultLang,
	}

	i, err = parameters.ParameterInt.Parser(ctx, 0, []string{"10"})
	require.Nil(t, err)
	require.Equal(t, i, 10)

	i, err = parameters.ParameterInt.Parser(ctx, 0, []string{"-10"})
	require.Nil(t, err)
	require.Equal(t, i, -10)

	i, err = parameters.ParameterInt.Parser(ctx, 1, []string{"hello", "10"})
	require.Nil(t, err)
	require.Equal(t, i, 10)

	_, err = parameters.ParameterInt.Parser(ctx, 0, []string{"hello"})
	require.NotNil(t, err)

	_, err = parameters.ParameterInt.Parser(ctx, 0, []string{"10.2"})
	require.NotNil(t, err)
}

func TestString(t *testing.T) {
	var str any
	var err error

	ctx := &command.Context{
		Lang: defaultLang,
	}

	str, err = parameters.ParameterString.Parser(ctx, 0, []string{"10"})
	require.Nil(t, err)
	require.Equal(t, str, "10")

	str, err = parameters.ParameterString.Parser(ctx, 1, []string{"hello", "-10"})
	require.Nil(t, err)
	require.Equal(t, str, "-10")
}

func TestText(t *testing.T) {
	var str any
	var err error

	ctx := &command.Context{
		Lang: defaultLang,
	}

	str, err = parameters.ParameterText.Parser(ctx, 0, []string{"hello", "world"})
	require.Nil(t, err)
	require.Equal(t, str, "hello world")

	str, err = parameters.ParameterText.Parser(ctx, 1, []string{"goodbye", "cruel", "world"})
	require.Nil(t, err)
	require.Equal(t, str, "cruel world")
}

func loadDefaultLang() {
	lang, err := i18n.GetLanguage(i18n.DefaultLanguageName)
	if err != nil {
		panic(err)
	}
	defaultLang = lang
}
