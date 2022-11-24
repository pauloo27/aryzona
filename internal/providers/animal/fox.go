package animal

import (
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/tidwall/gjson"
)

func GetRandomFox() (string, error) {
	json, err := utils.Get("https://randomfox.ca/floof/")
	if err != nil {
		return "", err
	}

	// rewrite with gjson
	return gjson.GetBytes(json, "image").String(), nil
}
