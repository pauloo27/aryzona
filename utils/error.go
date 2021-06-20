package utils

import "errors"

func HandleFatal(err error) {
	if err != nil {
		panic(err)
	}
}

func Wrap(msg string, err error) error {
	return errors.New(Fmt("%s: %v", msg, err))
}
