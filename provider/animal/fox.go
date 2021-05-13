package animal

import (
	"github.com/Pauloo27/aryzona/utils"
	"github.com/buger/jsonparser"
)

func GetRandomFox() (string, error) {
	json, err := utils.Get("https://randomfox.ca/floof/")
	if err != nil {
		return "", err
	}

	return jsonparser.GetString(json, "image")
}
