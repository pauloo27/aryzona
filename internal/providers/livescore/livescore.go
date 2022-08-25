package livescore

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pauloo27/logger"
	"github.com/buger/jsonparser"
)

// this code could use some refactoring...

type TeamInfo struct {
	Name, Color string
	Score       int
}

type Event struct {
	Text     string
	Min      int64
	ExtraMin int64
	Type     int64
}

func (t TeamInfo) ColorAsInt() int {
	color, err := strconv.ParseInt(t.Color, 16, 32)
	if err != nil {
		return 0
	}
	return int(color)
}

type MatchInfo struct {
	T1, T2                            *TeamInfo
	ID                                string
	CupName, StadiumName, StadiumCity string
	Time                              string // time as string? YES
	Events                            []*Event
}

func parseTeam(id int, data []byte) (*TeamInfo, error) {
	name, err := jsonparser.GetString(data, ("T" + strconv.Itoa(id)), "[0]", "Nm")
	if err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}

	color, err := jsonparser.GetString(data, ("T" + strconv.Itoa(id)), "[0]", "Shrt", "Bs")

	// sometimes the color is not defined... so lets use white as fallback
	if err != nil {
		color = "ffffff"
	}

	rawScore, _ := jsonparser.GetString(data, "Tr"+strconv.Itoa(id))

	score, err := strconv.Atoi(rawScore)
	if err != nil {
		score = 0
		//return nil, errore.Wrap("score", err)
	}

	return &TeamInfo{name, color, score}, nil
}

func parseMatchForListing(data []byte) (*MatchInfo, error) {
	id, err := jsonparser.GetString(data, "Eid")
	if err != nil {
		return nil, fmt.Errorf("id: %w", err)
	}

	time, err := jsonparser.GetString(data, "Eps")
	if err != nil {
		return nil, fmt.Errorf("time: %w", err)
	}

	team1, err := parseTeam(1, data)
	if err != nil {
		logger.Error("team1", err)
	}

	team2, err := parseTeam(2, data)
	if err != nil {
		return nil, fmt.Errorf("team2: %w", err)
	}

	return &MatchInfo{
		ID: id, T1: team1, T2: team2,
		Time: time,
	}, nil
}

func parseMatch(data []byte) (*MatchInfo, error) {
	id, err := jsonparser.GetString(data, "Eid")
	if err != nil {
		return nil, fmt.Errorf("id: %w", err)
	}

	time, err := jsonparser.GetString(data, "Eps")
	if err != nil {
		return nil, fmt.Errorf("time: %w", err)
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
		return nil, fmt.Errorf("team2: %w", err)
	}

	var events []*Event

	_, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		text, _ := jsonparser.GetString(value, "Txt")
		id, _ := jsonparser.GetInt(value, "IT")
		min, _ := jsonparser.GetInt(value, "Min")
		extraMin, _ := jsonparser.GetInt(value, "MinEx")

		event := Event{
			Text:     text,
			Type:     id,
			Min:      min,
			ExtraMin: extraMin,
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

	// TODO: refact with GJSON
	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		_, err0 := jsonparser.ArrayEach(value, func(matchData []byte, dataType jsonparser.ValueType, offset int, err error) {
			match, err1 := parseMatchForListing(matchData)
			if err1 != nil {
				logger.Fatal(err)
			} else {
				matches = append(matches, match)
			}
		}, "Events")
		if err0 != nil {
			logger.Error(err)
		}
	}, "Stages")
	if err != nil {
		return nil, fmt.Errorf("stages: %w", err)
	}

	return matches, nil
}
