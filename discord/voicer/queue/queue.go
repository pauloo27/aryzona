package queue

import (
	"github.com/Pauloo27/aryzona/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/utils/event"
)

const (
	EventAppend  = event.EventType("APPEND")
	EventPop     = event.EventType("POP")
	EventEnded   = event.EventType("ENDED")
	EventCleared = event.EventType("CLEARED")
)

type Queue struct {
	*event.EventEmitter
	queue []playable.Playable
}

func NewQueue() *Queue {
	return &Queue{
		queue:        []playable.Playable{},
		EventEmitter: event.NewEventEmitter(),
	}
}

func (q *Queue) Append(item playable.Playable) {
	q.queue = append(q.queue, item)
	q.Emit(EventAppend, q, q.Size()-1, item)
}

func (q *Queue) AppendAt(index int, item playable.Playable) {
	var tmp []playable.Playable
	tmp = append(tmp, q.queue[:index]...)
	tmp = append(tmp, item)
	tmp = append(tmp, q.queue[index:]...)
	q.queue = tmp
	q.Emit(EventAppend, q, index, item)
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

func (q *Queue) Pop(index int) {
	var tmp []playable.Playable
	tmp = append(tmp, q.queue[:index]...)
	tmp = append(tmp, q.queue[index+1:]...)
	q.queue = tmp
}

func (q *Queue) Size() int {
	return len(q.queue)
}
