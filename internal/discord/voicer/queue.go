package voicer

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/queue"
	"github.com/Pauloo27/aryzona/internal/utils/scheduler"
)

func (v *Voicer) IsPaused() bool {
	if v.StreamingSession == nil {
		return false
	}
	return v.StreamingSession.Paused()
}

func (v *Voicer) registerListeners() {
	start := func(params ...interface{}) {
		if v.IsPlaying() {
			return
		}
		scheduler.Unschedule(fmt.Sprintf("empty_queue_%s", *v.GuildID))
		_ = v.Start()
	}
	v.Queue.On(queue.EventAppend, start)

	v.Queue.On(queue.EventRemove, func(params ...interface{}) {
		data := params[0].(queue.EventRemoveData)
		if data.Index == 0 {
			v.EncodeSession.Cleanup()
		}
	})
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

func (v *Voicer) Playing() playable.Playable {
	return v.Queue.First()
}

func (v *Voicer) AppendToQueue(playable playable.Playable) error {
	v.Queue.Append(playable)
	return nil
}

func (v *Voicer) AppendManyToQueue(playable ...playable.Playable) error {
	v.Queue.AppendMany(playable...)
	return nil
}

func (v *Voicer) scheduleEmptyQueue() {
	task := scheduler.NewRunLaterTask(30*time.Second, func(params ...interface{}) {
		if v.IsConnected() {
			_ = v.Disconnect()
		}
	})
	scheduler.Schedule(fmt.Sprintf("empty_queue_%s", *v.GuildID), task)
}
