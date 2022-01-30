package parameters_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/stretchr/testify/require"
)

func TestBoolean(t *testing.T) {
	var b interface{}
	var err error

	b, err = parameters.ParameterBool.Parser(0, []string{"true"})
	require.Nil(t, err)
	require.Equal(t, b, true)

	b, err = parameters.ParameterBool.Parser(1, []string{"hello", "true"})
	require.Nil(t, err)
	require.Equal(t, b, true)

	b, err = parameters.ParameterBool.Parser(0, []string{"false"})
	require.Nil(t, err)
	require.Equal(t, b, false)

	b, err = parameters.ParameterBool.Parser(1, []string{"hello", "false"})
	require.Nil(t, err)
	require.Equal(t, b, false)

	_, err = parameters.ParameterBool.Parser(0, []string{"invalid"})
	require.NotNil(t, err)
}

func TestInt(t *testing.T) {
	var i interface{}
	var err error

	i, err = parameters.ParameterInt.Parser(0, []string{"10"})
	require.Nil(t, err)
	require.Equal(t, i, 10)

	i, err = parameters.ParameterInt.Parser(0, []string{"-10"})
	require.Nil(t, err)
	require.Equal(t, i, -10)

	i, err = parameters.ParameterInt.Parser(1, []string{"hello", "10"})
	require.Nil(t, err)
	require.Equal(t, i, 10)

	_, err = parameters.ParameterInt.Parser(0, []string{"hello"})
	require.NotNil(t, err)

	_, err = parameters.ParameterInt.Parser(0, []string{"10.2"})
	require.NotNil(t, err)
}

func TestString(t *testing.T) {
	var str interface{}
	var err error

	str, err = parameters.ParameterString.Parser(0, []string{"10"})
	require.Nil(t, err)
	require.Equal(t, str, "10")

	str, err = parameters.ParameterString.Parser(1, []string{"hello", "-10"})
	require.Nil(t, err)
	require.Equal(t, str, "-10")
}

func TestText(t *testing.T) {
	var str interface{}
	var err error

	str, err = parameters.ParameterText.Parser(0, []string{"hello", "world"})
	require.Nil(t, err)
	require.Equal(t, str, "hello world")

	str, err = parameters.ParameterText.Parser(1, []string{"goodbye", "cruel", "world"})
	require.Nil(t, err)
	require.Equal(t, str, "cruel world")
}
