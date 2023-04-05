package animal

import (
	"github.com/Pauloo27/aryzona/internal/core/h"
	"github.com/tidwall/gjson"
)

func GetRandomCat() (string, error) {
	json, err := h.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		return "", err
	}
	return gjson.GetBytes(json, "0.url").String(), nil
}
