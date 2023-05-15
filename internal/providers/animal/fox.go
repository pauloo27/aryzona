package animal

import (
	"github.com/pauloo27/aryzona/internal/core/h"
	"github.com/tidwall/gjson"
)

func GetRandomFox() (string, error) {
	json, err := h.Get("https://randomfox.ca/floof/")
	if err != nil {
		return "", err
	}

	// rewrite with gjson
	return gjson.GetBytes(json, "image").String(), nil
}
