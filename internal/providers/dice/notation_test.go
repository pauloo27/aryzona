package dice_test

import (
	"fmt"
	"testing"

	"github.com/pauloo27/aryzona/internal/providers/dice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// reference: https://en.wikipedia.org/wiki/Dice_notation
func parseAndCheck(t *testing.T, str string, ok bool, expectedDices, expectedFaces int) {
	should := "Should"
	if !ok {
		should += " not"
	}
	name := fmt.Sprintf("%s parse %s", should, str)

	t.Run(name, func(t *testing.T) {
		n, err := dice.ParseNotation(str)
		if !ok {
			assert.Nil(t, n)
			require.NotNil(t, err)
			return
		}
		assert.Nil(t, err)
		require.NotNil(t, n)

		assert.Equal(t, expectedDices, n.Dices)
		assert.Equal(t, expectedFaces, n.Faces)
	})
}

// eg: 10
func Test_RawNumber(t *testing.T) {
	parseAndCheck(t, "10", true, 1, 10)
	parseAndCheck(t, "32", true, 1, 32)
	parseAndCheck(t, "1", true, 1, 1)
	parseAndCheck(t, "0", false, 0, 0)
	parseAndCheck(t, "-10", false, 0, 0)
}

// eg: d10
func Test_dX(t *testing.T) {
	parseAndCheck(t, "d10", true, 1, 10)
	parseAndCheck(t, "d32", true, 1, 32)
	parseAndCheck(t, "d1", true, 1, 1)
	parseAndCheck(t, "d0", false, 0, 0)
	parseAndCheck(t, "d-10", false, 0, 0)
}

// eg: 2d10
func Test_AdX(t *testing.T) {
	parseAndCheck(t, "10d", true, 10, 6)
	parseAndCheck(t, "10d10", true, 10, 10)
	parseAndCheck(t, "32d8", true, 32, 8)
	parseAndCheck(t, "1d1", true, 1, 1)
	parseAndCheck(t, "1d6", true, 1, 6)
	parseAndCheck(t, "0d0", false, 0, 0)
	parseAndCheck(t, "0d-10", false, 0, 0)
	parseAndCheck(t, "-1d6", false, 0, 0)
	parseAndCheck(t, "-1d-1", false, 0, 0)
}

func Test_Extra(t *testing.T) {
	parseAndCheck(t, "d", true, 1, 6)
	parseAndCheck(t, "", false, 0, 0)
	parseAndCheck(t, "Ad2", false, 0, 0)
	parseAndCheck(t, "ad2", false, 0, 0)
	parseAndCheck(t, "2dA", false, 0, 0)
	parseAndCheck(t, "2da", false, 0, 0)
}
