package utils

import "errors"

type Errore struct {
	ID, Message string
}

func (v Errore) Error() string {
	return v.Message
}

func IsErrore(err error) (bool, Errore) {
	e, ok := err.(Errore)
	return ok, e
}

func Is(err error, vErr Errore) bool {
	voicerError, ok := err.(Errore)
	if !ok {
		return false
	}

	return voicerError.ID == vErr.ID
}

func HandleFatal(err error) {
	if err != nil {
		panic(err)
	}
}

func Wrap(msg string, err error) error {
	return errors.New(Fmt("%s: %v", msg, err))
}
