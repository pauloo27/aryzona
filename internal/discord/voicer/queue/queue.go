package queue

import (
	"math/rand"

	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/utils/event"
)

const (
	EventAppend  = event.EventType("APPEND")
	EventRemove  = event.EventType("REMOVE")
	EventShuffle = event.EventType("SHUFFLE")
)

type QueueEntry struct {
	Playable  playable.Playable
	Requester string
}

type Queue struct {
	queue []*QueueEntry
	*event.EventEmitter
}

func NewQueue() *Queue {
	return &Queue{
		queue:        []*QueueEntry{},
		EventEmitter: event.NewEventEmitter(),
	}
}

func (q *Queue) Shuffle() {
	for i := range q.queue {
		j := rand.Intn(i + 1) //#nosec G404
		// the top element on the list is the one being played,
		// there's no need to move it
		if i == 0 || j == 0 {
			continue
		}
		q.queue[i], q.queue[j] = q.queue[j], q.queue[i]
	}
	q.Emit(EventShuffle)
}

func (q *Queue) Append(item *QueueEntry) {
	q.queue = append(q.queue, item)
	q.Emit(EventAppend, EventAppendData{
		Queue:  q,
		Index:  q.Size() - 1,
		IsMany: false,
		Items:  []*QueueEntry{item},
	})
}

func (q *Queue) AppendMany(items ...*QueueEntry) {
	q.queue = append(q.queue, items...)
	q.Emit(EventAppend, EventAppendData{
		Queue:  q,
		Index:  q.Size() - 1,
		IsMany: true,
		Items:  items,
	})
}

func (q *Queue) All() []*QueueEntry {
	return q.queue
}

func (q *Queue) AppendAt(index int, item *QueueEntry) {
	var tmp []*QueueEntry
	tmp = append(tmp, q.queue[:index]...)
	tmp = append(tmp, item)
	tmp = append(tmp, q.queue[index:]...)
	q.queue = tmp
	q.Emit(EventAppend, EventAppendData{
		Queue:  q,
		Index:  index,
		IsMany: false,
		Items:  []*QueueEntry{item},
	})
}

func (q *Queue) Clear() {
	var tmp []*QueueEntry
	q.queue = tmp
}

func (q *Queue) ItemAt(index int) *QueueEntry {
	return q.queue[index]
}

func (q *Queue) First() *QueueEntry {
	if q.Size() == 0 {
		return nil
	}
	return q.queue[0]
}

func (q *Queue) Remove(index int) {
	if q.Size() == 0 {
		return
	}
	var tmp []*QueueEntry
	tmp = append(tmp, q.queue[:index]...)
	tmp = append(tmp, q.queue[index+1:]...)
	q.queue = tmp
	q.Emit(EventRemove, EventRemoveData{Queue: q, Index: index})
}

func (q *Queue) Size() int {
	return len(q.queue)
}
