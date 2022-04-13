package radio

import (
	"time"

	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
)

type RadioChannel interface {
	playable.Playable
	GetID() string
}

type BaseRadio struct {
}

func (c BaseRadio) CanPause() bool {
	return false
}

func (c BaseRadio) IsLive() bool {
	return true
}

func (BaseRadio) GetDuration() (time.Duration, error) {
	return 0, nil
}

func (c BaseRadio) IsLocal() bool {
	return false
}
