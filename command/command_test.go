package command

import (
	"testing"

	"github.com/Pauloo27/aryzona/utils"
	"github.com/stretchr/testify/assert"
)

func TestPermissions(t *testing.T) {
	t.Skip("No tests written yet...")
}

func TestArguments(t *testing.T) {
	t.Run("Required argument", func(t *testing.T) {
		testCommand := Command{
			Name: "Test command",
			Arguments: []*CommandArgument{
				{
					Name:     "test string",
					Required: true,
					Type:     ArgumentString,
				},
			},
		}

		t.Run("Should return required argument missing", func(t *testing.T) {
			_, err := testCommand.ValidateArguments([]string{})
			assert.NotNil(t, err)
			assert.True(t, utils.Is(*err, *ErrRequiredArgument(nil)))
		})

		t.Run("Should return the value", func(t *testing.T) {
			input := "hello"
			values, err := testCommand.ValidateArguments([]string{input})
			assert.Nil(t, err)
			assert.NotNil(t, values)
			assert.Len(t, values, 1)
			assert.Equal(t, input, values[0])
		})
	})

	t.Run("Not required argument", func(t *testing.T) {
		testCommand := Command{
			Name: "Test command",
			Arguments: []*CommandArgument{
				{
					Name:     "test string",
					Required: false,
					Type:     ArgumentString,
				},
			},
		}

		t.Run("Should NOT return required argument missing", func(t *testing.T) {
			values, err := testCommand.ValidateArguments([]string{})
			assert.Nil(t, err)
			assert.Nil(t, values)
			assert.Len(t, values, 0)
		})

		t.Run("Should return the value", func(t *testing.T) {
			input := "hello"
			values, err := testCommand.ValidateArguments([]string{input})
			assert.Nil(t, err)
			assert.NotNil(t, values)
			assert.Len(t, values, 1)
			assert.Equal(t, input, values[0])
		})
	})

	t.Run("Parse int argument", func(t *testing.T) {
		testCommand := Command{
			Name: "Test command",
			Arguments: []*CommandArgument{
				{
					Name:     "test string",
					Required: false,
					Type:     ArgumentInt,
				},
			},
		}

		t.Run("Should return cannot parse argument", func(t *testing.T) {
			values, err := testCommand.ValidateArguments([]string{"asd"})
			assert.NotNil(t, err)
			assert.Nil(t, values)
			assert.True(t, utils.Is(*err, *ErrCannotParseArgument(nil, nil)))
		})
	})

	// TODO: test with "valid values"
	t.Run("Check if value is inside a list", func(t *testing.T) {
		t.Skip("Not implemeneted yet")
	})

	// TODO: complex command
}
