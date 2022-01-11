package livescore

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/internal/utils/errore"
	"github.com/Pauloo27/logger"
	"github.com/buger/jsonparser"
)

type TeamInfo struct {
	Name, Color string
	Score       int
}

type Event struct {
	Text string
	Min  int64
	ID   int64
}

func (t TeamInfo) ColorAsInt() int {
	color, err := strconv.ParseInt(t.Color, 16, 32)
	if err != nil {
		return 0
	}
	return int(color)
}

type MatchInfo struct {
	Events                            []*Event
	ID                                string
	Time                              string // time as string? YES
	CupName, StadiumName, StadiumCity string
	T1, T2                            *TeamInfo
}

func parseTeam(id int, data []byte) (*TeamInfo, error) {
	name, err := jsonparser.GetString(data, ("T" + strconv.Itoa(id)), "[0]", "Nm")
	if err != nil {
		return nil, errore.Wrap("name", err)
	}

	color, err := jsonparser.GetString(data, ("T" + strconv.Itoa(id)), "[0]", "Shrt", "Bs")

	// sometimes the color is not defined... so lets use white as fallback
	if err != nil {
		color = "ffffff"
	}

	rawScore, err := jsonparser.GetString(data, "Tr"+strconv.Itoa(id))
	if err != nil {
		logger.Warn(errore.Wrap("cannot get raw score", err))
	}

	score, err := strconv.Atoi(rawScore)
	if err != nil {
		score = -1
		//return nil, errore.Wrap("score", err)
	}

	return &TeamInfo{name, color, score}, nil
}

func parseMatchForListing(data []byte) (*MatchInfo, error) {
	id, err := jsonparser.GetString(data, "Eid")
	if err != nil {
		return nil, errore.Wrap("id", err)
	}

	team1, err := parseTeam(1, data)
	if err != nil {
		logger.Error("team1", err)
	}

	team2, err := parseTeam(2, data)
	if err != nil {
		return nil, errore.Wrap("team2", err)
	}

	return &MatchInfo{
		ID: id, T1: team1, T2: team2,
	}, nil
}

func parseMatch(data []byte) (*MatchInfo, error) {
	id, err := jsonparser.GetString(data, "Eid")
	if err != nil {
		return nil, errore.Wrap("id", err)
	}

	time, err := jsonparser.GetString(data, "Eps")
	if err != nil {
		return nil, errore.Wrap("time", err)
	}

	cupName, err := jsonparser.GetString(data, "Stg", "Sdn")
	if err != nil {
		cupName = "Stree game"
	}

	stadiumName, err := jsonparser.GetString(data, "Vnm")
	if err != nil {
		stadiumName = "Somewhere"
	}

	stadiumCity, err := jsonparser.GetString(data, "VCity")
	if err != nil {
		stadiumCity = "Earth"
	}

	team1, err := parseTeam(1, data)
	if err != nil {
		logger.Error("team1", err)
	}

	team2, err := parseTeam(2, data)
	if err != nil {
		return nil, errore.Wrap("team2", err)
	}

	var events []*Event

	_, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		text, _ := jsonparser.GetString(value, "Txt")
		id, _ := jsonparser.GetInt(value, "IT")
		min, _ := jsonparser.GetInt(value, "Min")

		event := Event{
			Text: text,
			ID:   id,
			Min:  min,
		}
		events = append(events, &event)
	}, "Com")

	return &MatchInfo{
		ID: id, T1: team1, T2: team2, Time: time, CupName: cupName,
		StadiumName: stadiumName, StadiumCity: stadiumCity,
		Events: events,
	}, nil
}

func FetchMatchInfoByTeamName(teamName string) (*MatchInfo, error) {
	matches, err := ListLives()
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		if strings.EqualFold(match.T1.Name, teamName) ||
			strings.EqualFold(match.T2.Name, teamName) {
			return FetchMatchInfo(match.ID)
		}
	}
	return nil, nil
}

/* #nosec GG107 */
func FetchMatchInfo(matchID string) (*MatchInfo, error) {
	endpoint := fmt.Sprintf("https://prod-public-api.livescore.com/v1/api/react/match-x/soccer/%s/-3", matchID)

	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseMatch(data)
}

func ListLives() ([]*MatchInfo, error) {
	endpoint := "https://prod-public-api.livescore.com/v1/api/react/live/soccer/-3.00"

	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	matches := []*MatchInfo{}

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		_, err0 := jsonparser.ArrayEach(value, func(matchData []byte, dataType jsonparser.ValueType, offset int, err error) {
			match, err1 := parseMatchForListing(matchData)
			if err1 != nil {
				errore.HandleFatal(err)
			} else {
				matches = append(matches, match)
			}
		}, "Events")
		if err0 != nil {
			logger.Error(err)
		}
	}, "Stages")
	if err != nil {
		return nil, errore.Wrap("stages", err)
	}

	return matches, nil
}
