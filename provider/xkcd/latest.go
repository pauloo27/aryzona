package xkcd

import (
	"github.com/Pauloo27/aryzona/utils"
	"github.com/buger/jsonparser"
)

func GetLatest() (string, error) {
	json, err := utils.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return "", err
	}

	return jsonparser.GetString(json, "img")
}
