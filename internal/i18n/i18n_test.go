package i18n_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/stretchr/testify/assert"
)

func TestEntryFormat(t *testing.T) {
	var entry i18n.Entry

	entry = i18n.Entry("Hello")
	assert.Equal(t, "Hello", entry.Str())

	entry = i18n.Entry("Hello {0:name}")
	assert.Equal(t, "Hello {0:name}", entry.Str())

	entry = i18n.Entry("Hello {0:name}")
	assert.Equal(t, "Hello World", entry.Str("World"))

	entry = i18n.Entry("Hello {0:name}")
	assert.Equal(t, "Hello World", entry.Str("World", "Invalid"))

	entry = i18n.Entry("Hello {0:name}, {1:greet}")
	assert.Equal(t, "Hello World, welcome", entry.Str("World", "welcome"))

	entry = i18n.Entry("Hello {0:name}, {1:greet}")
	assert.Equal(t, "Hello World, welcome", entry.Str("World", "welcome"))

	entry = i18n.Entry("Hello {0:name}, {1:greet}")
	assert.Equal(t, "Hello 123, true", entry.Str(123, true))
}
