package voicer

import (
	"github.com/Pauloo27/aryzona/internal/utils/errore"
)

var (
	ErrAlreadyPlaying = errore.Errore{
		ID:      "ALREADY_PLAYING",
		Message: "Already playing something in the current guild",
	}
	ErrCannotConnect = errore.Errore{
		ID:      "CANNOT_CONNECT",
		Message: "Cannot connect to voice channel",
	}
)
