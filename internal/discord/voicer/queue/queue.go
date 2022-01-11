package queue

import (
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/utils/event"
)

const (
	EventAppend = event.EventType("APPEND")
	EventRemove = event.EventType("REMOVE")
)

type Queue struct {
	queue []playable.Playable
	*event.EventEmitter
}

func NewQueue() *Queue {
	return &Queue{
		queue:        []playable.Playable{},
		EventEmitter: event.NewEventEmitter(),
	}
}

func (q *Queue) Append(item playable.Playable) {
	q.queue = append(q.queue, item)
	q.Emit(EventAppend, EventAppendData{
		Queue:  q,
		Index:  q.Size() - 1,
		IsMany: false,
		Items:  []playable.Playable{item},
	})
}

func (q *Queue) AppendMany(items ...playable.Playable) {
	q.queue = append(q.queue, items...)
	q.Emit(EventAppend, EventAppendData{
		Queue:  q,
		Index:  q.Size() - 1,
		IsMany: true,
		Items:  items,
	})
}

func (q *Queue) All() []playable.Playable {
	return q.queue
}

func (q *Queue) AppendAt(index int, item playable.Playable) {
	var tmp []playable.Playable
	tmp = append(tmp, q.queue[:index]...)
	tmp = append(tmp, item)
	tmp = append(tmp, q.queue[index:]...)
	q.queue = tmp
	q.Emit(EventAppend, EventAppendData{
		Queue:  q,
		Index:  index,
		IsMany: false,
		Items:  []playable.Playable{item},
	})
}

func (q *Queue) Clear() {
	var tmp []playable.Playable
	q.queue = tmp
}

func (q *Queue) ItemAt(index int) playable.Playable {
	return q.queue[index]
}

func (q *Queue) First() playable.Playable {
	if q.Size() == 0 {
		return nil
	}
	return q.queue[0]
}

func (q *Queue) Remove(index int) {
	if q.Size() == 0 {
		return
	}
	var tmp []playable.Playable
	tmp = append(tmp, q.queue[:index]...)
	tmp = append(tmp, q.queue[index+1:]...)
	q.queue = tmp
	q.Emit(EventRemove, EventRemoveData{Queue: q, Index: index})
}

func (q *Queue) Size() int {
	return len(q.queue)
}
