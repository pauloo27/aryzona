package queue_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/queue"
	"github.com/stretchr/testify/assert"
)

// yeap, tests for Append and Remove operations =P
func TestQueue(t *testing.T) {
	q := queue.NewQueue()
	assert.Equal(t, 0, q.Size())
	assert.Nil(t, q.First())

	q.Append(
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "hello"},
			Requester: "Dummy",
		},
	)
	q.Append(
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "coming next"},
			Requester: "Dummy",
		},
	)
	assert.Equal(t, 2, q.Size())
	assert.NotNil(t, q.First())
	assert.Equal(t, "hello", q.First().Playable.GetName())

	q.Remove(0)
	assert.Equal(t, 1, q.Size())
	assert.NotNil(t, q.First())
	assert.Equal(t, "coming next", q.First().Playable.GetName())

	q.Append(
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "hello"},
			Requester: "Dummy",
		},
	)
	q.Append(
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "bye"},
			Requester: "Dummy",
		},
	)

	assert.Equal(t, 3, q.Size())
	assert.NotNil(t, q.ItemAt(1))
	assert.Equal(t, "hello", q.ItemAt(1).Playable.GetName())

	q.Remove(1)
	assert.Equal(t, 2, q.Size())
	assert.NotNil(t, q.First())
	assert.Equal(t, "coming next", q.First().Playable.GetName())

	q.Remove(0)
	assert.Equal(t, 1, q.Size())
	assert.NotNil(t, q.First())
	assert.Equal(t, "bye", q.First().Playable.GetName())

	q.Clear()
	assert.Equal(t, 0, q.Size())
	assert.Nil(t, q.First())

	q.Append(
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "hello"},
			Requester: "Dummy",
		},
	)
	q.Append(
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "bye"},
			Requester: "Dummy",
		},
	)
	q.AppendAt(1, &queue.QueueEntry{
		Playable:  playable.DummyPlayable{Name: "=)"},
		Requester: "Dummy",
	})

	assert.Equal(t, 3, q.Size())
	assert.Equal(t, "hello", q.ItemAt(0).Playable.GetName())
	assert.Equal(t, "=)", q.ItemAt(1).Playable.GetName())
	assert.Equal(t, "bye", q.ItemAt(2).Playable.GetName())

	assert.Equal(t, 3, len(q.All()))
	assert.Equal(t, "hello", q.All()[0].Playable.GetName())
	assert.Equal(t, "=)", q.All()[1].Playable.GetName())
	assert.Equal(t, "bye", q.All()[2].Playable.GetName())

	q.Clear()
	assert.Equal(t, 0, q.Size())

	q.AppendMany(
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "hello"},
			Requester: "Dummy",
		},
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "bye"},
			Requester: "Dummy",
		},
		&queue.QueueEntry{
			Playable:  playable.DummyPlayable{Name: "=)"},
			Requester: "Dummy",
		},
	)

	assert.Equal(t, 3, q.Size())
	assert.Equal(t, "hello", q.All()[0].Playable.GetName())
	assert.Equal(t, "bye", q.All()[1].Playable.GetName())
	assert.Equal(t, "=)", q.All()[2].Playable.GetName())
}

// TODO: test events
