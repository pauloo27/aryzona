package xkcd

import (
	"encoding/json"
	"fmt"

	"github.com/Pauloo27/aryzona/internal/utils"
)

type Comic struct {
	Num        int    `json:"num"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Alt        string `json:"alt"`
	Year       string `json:"year"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	Transcript string `json:"transcript"`
	SafeTitle  string `json:"safe_title"`
	News       string `json:"news"`
}

func GetByNum(num int) (*Comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num)
	if num == 0 {
		url = "https://xkcd.com/info.0.json"
	}
	res, err := utils.Get(url)
	if err != nil {
		return nil, err
	}

	var xkcd Comic

	err = json.Unmarshal([]byte(res), &xkcd)
	if err != nil {
		return nil, err
	}
	return &xkcd, nil
}
