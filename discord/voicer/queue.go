package voicer

import "github.com/Pauloo27/aryzona/audio"

type Queue struct {
	queue []audio.Playable
}

func (q *Queue) Append(item audio.Playable) {
	q.queue = append(q.queue, item)
}

func (q *Queue) Playing() audio.Playable {
	if q.Size() == 0 {
		return nil
	}
	return q.queue[0]
}

func (q *Queue) Pop(index int) {
	var tmp []audio.Playable
	tmp = q.queue[:index]
	tmp = append(tmp, q.queue[index+1:]...)
	q.queue = tmp
}

func (q *Queue) Size() int {
	return len(q.queue)
}
