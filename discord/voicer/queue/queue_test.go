package queue_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/discord/voicer/queue"
	"github.com/stretchr/testify/assert"
)

// yeap, a test for Append and Pop operations =P
func TestQueue(t *testing.T) {
	queue := queue.NewQueue()
	assert.Equal(t, 0, queue.Size())
	assert.Nil(t, queue.First())

	queue.Append(playable.DummyPlayable{Name: "hello"})
	queue.Append(playable.DummyPlayable{Name: "coming next"})
	assert.Equal(t, 2, queue.Size())
	assert.NotNil(t, queue.First())
	assert.Equal(t, "hello", queue.First().GetName())

	queue.Remove(0)
	assert.Equal(t, 1, queue.Size())
	assert.NotNil(t, queue.First())
	assert.Equal(t, "coming next", queue.First().GetName())

	queue.Append(playable.DummyPlayable{Name: "hello"})
	queue.Append(playable.DummyPlayable{Name: "bye"})

	assert.Equal(t, 3, queue.Size())
	assert.NotNil(t, queue.ItemAt(1))
	assert.Equal(t, "hello", queue.ItemAt(1).GetName())

	queue.Remove(1)
	assert.Equal(t, 2, queue.Size())
	assert.NotNil(t, queue.First())
	assert.Equal(t, "coming next", queue.First().GetName())

	queue.Remove(0)
	assert.Equal(t, 1, queue.Size())
	assert.NotNil(t, queue.First())
	assert.Equal(t, "bye", queue.First().GetName())

	queue.Clear()
	assert.Equal(t, 0, queue.Size())
	assert.Nil(t, queue.First())

	queue.Append(playable.DummyPlayable{Name: "hello"})
	queue.Append(playable.DummyPlayable{Name: "bye"})
	queue.AppendAt(1, playable.DummyPlayable{Name: "=)"})

	assert.Equal(t, 3, queue.Size())
	assert.Equal(t, "hello", queue.ItemAt(0).GetName())
	assert.Equal(t, "=)", queue.ItemAt(1).GetName())
	assert.Equal(t, "bye", queue.ItemAt(2).GetName())

	assert.Equal(t, 3, len(queue.All()))
	assert.Equal(t, "hello", queue.All()[0].GetName())
	assert.Equal(t, "=)", queue.All()[1].GetName())
	assert.Equal(t, "bye", queue.All()[2].GetName())
}
