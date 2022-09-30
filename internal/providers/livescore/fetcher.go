package livescore

import (
	"errors"
	"sync"
)

type LiveMatch struct {
	MatchID      string
	CurrentData  MatchInfo
	PreviousData MatchInfo
}

var (
	ErrorMatchAlreadyFollowed = errors.New("match already followed")
)

var (
	FollowedMatches = make(map[string]*LiveMatch)
	mutex           = &sync.Mutex{}
)

func FollowMatch(id string) (*LiveMatch, error) {
	mutex.Lock()
	defer mutex.Unlock()

	m, found := FollowedMatches[id]
	if found {
		return m, ErrorMatchAlreadyFollowed
	}
	// TODO:
	return nil, nil
}
