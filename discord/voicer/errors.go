package voicer

import "github.com/Pauloo27/aryzona/utils"

var (
	ErrAlreadyPlaying = utils.Errore{
		ID:      "ALREADY_PLAYING",
		Message: "Already playing something in the current guild",
	}
	ErrCannotConnect = utils.Errore{
		ID:      "CANNOT_CONNECT",
		Message: "Cannot connect to voice channel",
	}
)
