package voicer

import (
	"fmt"

	"github.com/Pauloo27/aryzona/audio"
)

type Queue struct {
	queue []audio.Playable
}

func (q *Queue) Append(item audio.Playable) {
	q.queue = append(q.queue, item)
}

func (q *Queue) AppendAfter(index int, item audio.Playable) {
	var tmp []audio.Playable
	tmp = append(tmp, q.queue[:index+1]...)
	tmp = append(tmp, item)
	tmp = append(tmp, q.queue[index+1:]...)
	fmt.Println(tmp)
	q.queue = tmp
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
