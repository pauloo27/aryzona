package voicer_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/stretchr/testify/assert"
)

// yeap, a test for Append and Pop operations =P
func TestQueue(t *testing.T) {
	queue := voicer.Queue{}
	assert.Equal(t, 0, queue.Size())
	assert.Nil(t, queue.Playing())

	queue.Append(TestPlayable{Name: "hello"})
	queue.Append(TestPlayable{Name: "coming next"})
	assert.Equal(t, 2, queue.Size())
	assert.NotNil(t, queue.Playing())
	assert.Equal(t, "hello", queue.Playing().GetName())

	queue.Pop(0)
	assert.Equal(t, 1, queue.Size())
	assert.NotNil(t, queue.Playing())
	assert.Equal(t, "coming next", queue.Playing().GetName())

	queue.Append(TestPlayable{Name: "hello"})
	queue.Append(TestPlayable{Name: "bye"})

	assert.Equal(t, 3, queue.Size())
	assert.NotNil(t, queue.ItemAt(1))
	assert.Equal(t, "hello", queue.ItemAt(1).GetName())

	queue.Pop(1)
	assert.Equal(t, 2, queue.Size())
	assert.NotNil(t, queue.Playing())
	assert.Equal(t, "coming next", queue.Playing().GetName())

	queue.Pop(0)
	assert.Equal(t, 1, queue.Size())
	assert.NotNil(t, queue.Playing())
	assert.Equal(t, "bye", queue.Playing().GetName())

	queue.Clear()
	assert.Equal(t, 0, queue.Size())
	assert.Nil(t, queue.Playing())

	queue.Append(TestPlayable{Name: "hello"})
	queue.Append(TestPlayable{Name: "bye"})
	queue.AppendAfter(0, TestPlayable{Name: "=)"})

	assert.Equal(t, 3, queue.Size())
	assert.Equal(t, "hello", queue.ItemAt(0).GetName())
	assert.Equal(t, "=)", queue.ItemAt(1).GetName())
	assert.Equal(t, "bye", queue.ItemAt(2).GetName())
}
