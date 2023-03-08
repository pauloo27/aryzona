package voicer

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/core/scheduler"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/queue"
)

func (v *Voicer) IsPaused() bool {
	if v.StreamingSession == nil {
		return false
	}
	return v.StreamingSession.Paused()
}

func (v *Voicer) registerListeners() {
	start := func(params ...any) {
		if v.IsPlaying() {
			return
		}
		scheduler.Unschedule(fmt.Sprintf("empty_queue_%s", *v.GuildID))
		_ = v.Start()
	}
	v.Queue.On(queue.EventAppend, start)

	v.Queue.On(queue.EventRemove, func(params ...any) {
		data := params[0].(queue.EventRemoveData)
		if data.Index == 0 {
			v.EncodeSession.Cleanup()
		}
	})
}

func (v *Voicer) Pause() {
	if v.StreamingSession == nil {
		return
	}
	v.StreamingSession.Pause()
}

func (v *Voicer) Resume() {
	if v.StreamingSession == nil {
		return
	}
	v.StreamingSession.Resume()
}

func (v *Voicer) TogglePause() {
	if v.StreamingSession == nil {
		return
	}
	v.StreamingSession.TogglePause()
}

func (v *Voicer) Skip() {
	if v.EncodeSession == nil {
		return
	}
	v.EncodeSession.Cleanup()
}

func (v *Voicer) Playing() *queue.QueueEntry {
	return v.Queue.First()
}

func (v *Voicer) AppendToQueue(requesterID string, playable playable.Playable) error {
	v.Queue.Append(&queue.QueueEntry{
		Playable:  playable,
		Requester: requesterID,
	})
	return nil
}

func (v *Voicer) AppendManyToQueue(requesterID string, playable ...playable.Playable) error {
	var entries []*queue.QueueEntry
	for _, p := range playable {
		entries = append(entries, &queue.QueueEntry{
			Playable:  p,
			Requester: requesterID,
		})
	}
	v.Queue.AppendMany(entries...)
	return nil
}

func (v *Voicer) scheduleEmptyQueue() {
	task := scheduler.NewRunLaterTask(30*time.Second, func(params ...any) {
		if v.IsConnected() {
			_ = v.Disconnect()
		}
	})
	scheduler.Schedule(fmt.Sprintf("empty_queue_%s", *v.GuildID), task)
}
