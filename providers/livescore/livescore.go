package livescore

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Pauloo27/aryzona/utils"
	"github.com/buger/jsonparser"
)

type TeamInfo struct {
	Name, Color string
	Score       int
}

func (t TeamInfo) ColorAsInt() int {
	color, err := strconv.ParseInt(t.Color, 16, 32)
	if err != nil {
		return 0
	}
	return int(color)
}

type MatchInfo struct {
	Id                                string
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
		//return nil, utils.Wrap("color", err)
		color = "ffffff"
	}

	rawScore, err := jsonparser.GetString(data, "Tr"+strconv.Itoa(id))
	if err != nil {
		//return nil, utils.Wrap("raw score", err)
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
		//return nil, utils.Wrap("stadium name", err)
	}

	stadiumCity, err := jsonparser.GetString(data, "VCity")
	if err != nil {
		//return nil, utils.Wrap("stadium city", err)
	}

	team1, err := parseTeam(1, data)
	if err != nil {
		//return nil, utils.Wrap("team1", err)
	}

	team2, err := parseTeam(2, data)
	if err != nil {
		return nil, utils.Wrap("team2", err)
	}

	return &MatchInfo{
		id, team1, team2, time, cupName, stadiumName, stadiumCity,
	}, nil
}

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
		jsonparser.ArrayEach(value, func(matchData []byte, dataType jsonparser.ValueType, offset int, err error) {
			match, err := parseMatch(matchData)
			if err != nil {
				// TODO???? WHATTTTTT
				utils.HandleFatal(err)
			} else {
				matches = append(matches, match)
			}
		}, "Events")
	}, "Stages")
	if err != nil {
		return nil, utils.Wrap("stages", err)
	}

	return matches, nil
}