package voicer

import (
	"fmt"

	"github.com/Pauloo27/aryzona/audio"
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
	queue []audio.Playable
}

func NewQueue() *Queue {
	return &Queue{
		queue:        []audio.Playable{},
		EventEmitter: event.NewEventEmitter(),
	}
}

func (q *Queue) Append(item audio.Playable) {
	q.queue = append(q.queue, item)
	q.Emit(EventAppend, q, 0, item)
}

func (q *Queue) AppendAt(index int, item audio.Playable) {
	var tmp []audio.Playable
	tmp = append(tmp, q.queue[:index]...)
	tmp = append(tmp, item)
	tmp = append(tmp, q.queue[index:]...)
	fmt.Println(tmp)
	q.queue = tmp
	q.Emit(EventAppend, q, index, item)
}

func (q *Queue) Clear() {
	var tmp []audio.Playable
	q.queue = tmp
}

func (q *Queue) ItemAt(index int) audio.Playable {
	return q.queue[index]
}

func (q *Queue) First() audio.Playable {
	if q.Size() == 0 {
		return nil
	}
	return q.queue[0]
}

func (q *Queue) Pop(index int) {
	var tmp []audio.Playable
	tmp = append(tmp, q.queue[:index]...)
	tmp = append(tmp, q.queue[index+1:]...)
	q.queue = tmp
}

func (q *Queue) Size() int {
	return len(q.queue)
}
