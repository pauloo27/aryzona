package animal

import (
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/buger/jsonparser"
)

func GetRandomCat() (string, error) {
	json, err := utils.Get("https://aws.random.cat/meow")
	if err != nil {
		return "", err
	}

	return jsonparser.GetString(json, "file")
}
