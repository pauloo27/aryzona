package xkcd

import (
	"math/rand"
)

func GetRandom() (*Comic, error) {
	latest, err := GetLatest()
	if err != nil {
		return nil, err
	}
	/* #nosec G404 */
	return GetByNum(rand.Intn(latest.Num))
}
