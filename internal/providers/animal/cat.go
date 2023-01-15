package animal

import (
	"github.com/Pauloo27/aryzona/internal/core/h"
	"github.com/tidwall/gjson"
)

func GetRandomCat() (string, error) {
	json, err := h.Get("https://aws.random.cat/meow")
	if err != nil {
		return "", err
	}

	return gjson.GetBytes(json, "file").String(), nil
}
