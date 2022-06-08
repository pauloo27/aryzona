package rnd

import (
	"crypto/rand"
	"math/big"
)

func Rnd(n int) (int, error) {
	bigLuckyNumber, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}
	return int(bigLuckyNumber.Int64()), nil
}
