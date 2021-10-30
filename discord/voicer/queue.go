package voicer

import "github.com/Pauloo27/aryzona/discord/voicer/playable"

func (v *Voicer) IsPaused() bool {
	if v.StreamingSession == nil {
		return false
	}
	return v.StreamingSession.Paused()
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
