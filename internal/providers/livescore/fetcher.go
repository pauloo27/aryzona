package livescore

import (
	"errors"
	"time"

	"github.com/Pauloo27/aryzona/internal/utils"
)

type ListenerFn func(*LiveMatch, error)

type LiveMatch struct {
	MatchID      string
	CurrentData  *MatchInfo
	PreviousData *MatchInfo
	Listeners    []ListenerFn
}

// errors
var (
	ErrMatchAlreadyFollowed = errors.New("match already followed")
	ErrMatchHasFinished     = errors.New("match has finished")
)

var (
	followedMatches = make(map[string]*LiveMatch)
)

const (
	updatePeriod = 30 * time.Second
)

func init() {
	utils.Go(func() {
		for {
			for _, match := range followedMatches {
				updateLiveMatch(match)
			}
			time.Sleep(updatePeriod)
		}
	})
}

func GetLiveMatch(id string) (*LiveMatch, error) {
	match, found := followedMatches[id]
	if found {
		return match, nil
	}

	return followMatch(id)
}

func followMatch(id string) (*LiveMatch, error) {
	match, err := FetchMatchInfo(id)
	if err != nil {
		return nil, err
	}

	if match.Time == "FT" {
		return nil, ErrMatchHasFinished
	}

	followedMatches[id] = &LiveMatch{
		MatchID:      id,
		CurrentData:  match,
		PreviousData: nil,
	}

	return followedMatches[id], nil
}

func UnfollowMatch(id string) {
	delete(followedMatches, id)
}

func updateLiveMatch(liveMatch *LiveMatch) {
	match, err := FetchMatchInfo(liveMatch.MatchID)
	if err != nil {
		for _, listener := range liveMatch.Listeners {
			listener(nil, err)
		}
		UnfollowMatch(liveMatch.MatchID)
	}

	liveMatch.PreviousData = liveMatch.CurrentData
	liveMatch.CurrentData = match

	for _, listener := range liveMatch.Listeners {
		listener(liveMatch, nil)
	}

	if match.Time == "FT" {
		UnfollowMatch(liveMatch.MatchID)
		for _, listener := range liveMatch.Listeners {
			listener(liveMatch, ErrMatchHasFinished)
		}
	}
}

func (liveMatch *LiveMatch) AddListener(listener ListenerFn) {
	liveMatch.Listeners = append(liveMatch.Listeners, listener)
}
