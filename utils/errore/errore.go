package errore

import (
	"fmt"
)

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
	errore, ok := err.(Errore)
	if !ok {
		return false
	}

	return errore.ID == vErr.ID
}

func HandleFatal(err error) {
	if err != nil {
		panic(err)
	}
}

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
