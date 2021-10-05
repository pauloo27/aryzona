package livescore

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
	"github.com/buger/jsonparser"
)

type TeamInfo struct {
	Name, Color string
	Score       int
}

type Event struct {
	Min  int64
	Text string
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
	T1, T2                            *TeamInfo
	Time                              string // time as string? YES
	CupName, StadiumName, StadiumCity string
}

func parseTeam(id int, data []byte) (*TeamInfo, error) {
	name, err := jsonparser.GetString(data, ("T" + strconv.Itoa(id)), "[0]", "Nm")
	if err != nil {
		return nil, utils.Wrap("name", err)
	}

	color, err := jsonparser.GetString(data, ("T" + strconv.Itoa(id)), "[0]", "Shrt", "Bs")

	// sometimes the color is not defined... so lets use white as fallback
	if err != nil {
		color = "ffffff"
		logger.Error(err)
	}

	rawScore, err := jsonparser.GetString(data, "Tr"+strconv.Itoa(id))
	if err != nil {
		logger.Error(err)
	}

	score, err := strconv.Atoi(rawScore)
	if err != nil {
		score = -1
		//return nil, utils.Wrap("score", err)
	}

	return &TeamInfo{name, color, score}, nil
}

func parseMatch(data []byte) (*MatchInfo, error) {
	id, err := jsonparser.GetString(data, "Eid")
	if err != nil {
		return nil, utils.Wrap("id", err)
	}

	time, err := jsonparser.GetString(data, "Eps")
	if err != nil {
		return nil, utils.Wrap("time", err)
	}

	cupName, err := jsonparser.GetString(data, "Stg", "Sdn")
	if err != nil {
		return nil, utils.Wrap("cup name", err)
	}

	stadiumName, err := jsonparser.GetString(data, "Vnm")
	if err != nil {
		logger.Error("stadium name", err)
	}

	stadiumCity, err := jsonparser.GetString(data, "VCity")
	if err != nil {
		logger.Error("stadium city", err)
	}

	team1, err := parseTeam(1, data)
	if err != nil {
		logger.Error("team1", err)
	}

	team2, err := parseTeam(2, data)
	if err != nil {
		return nil, utils.Wrap("team2", err)
	}

	var events []*Event

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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

	if err != nil {
		return nil, err
	}

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
			return match, nil
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
			match, err1 := parseMatch(matchData)
			if err1 != nil {
				utils.HandleFatal(err)
			} else {
				matches = append(matches, match)
			}
		}, "Events")
		if err0 != nil {
			logger.Error(err)
		}
	}, "Stages")
	if err != nil {
		return nil, utils.Wrap("stages", err)
	}

	return matches, nil
}
