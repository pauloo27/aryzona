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
	// TODO: test required
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
	// TODO: test not required
	t.Run("Not required argument", func(t *testing.T) {
		t.Fatal("No tests written yet")
	})
	// TODO: test parser
	t.Run("Parse int argument", func(t *testing.T) {
		t.Fatal("No tests written yet")
	})
	// TODO: test with "valid values"
	t.Run("Check if value is inside a list", func(t *testing.T) {
		t.Fatal("No tests written yet")
	})
}
