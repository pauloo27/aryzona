package voicer

import (
	"time"

	"github.com/Pauloo27/aryzona/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/discord/voicer/queue"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/aryzona/utils/scheduler"
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
		scheduler.Unschedule(utils.Fmt("empty_queue_%s", *v.GuildID))
		_ = v.Start()
	}
	v.Queue.On(queue.EventAppend, start)

	v.Queue.On(queue.EventPop, func(params ...interface{}) {
		index := params[1].(int)
		if index == 0 {
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

func (v *Voicer) scheduleEmptyQueue() {
	task := scheduler.NewRunLaterTask(30*time.Second, func(params ...interface{}) {
		if v.IsConnected() {
			_ = v.Disconnect()
		}
	})
	scheduler.Schedule(utils.Fmt("empty_queue_%s", *v.GuildID), task)
}