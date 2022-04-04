package xkcd

import (
	"math/rand"
	"time"
)

func GetRandom() (*Comic, error) {
	latest, err := GetLatest()
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().UnixNano())
	/* #nosec G404 */
	return GetByNum(rand.Intn(latest.Num))
}
