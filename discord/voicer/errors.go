package voicer

import "github.com/Pauloo27/aryzona/utils"

var (
	ERR_ALREADY_PLAYING = utils.Errore{
		ID:      "ALREADY_PLAYING",
		Message: "Already playing something in the current guild",
	}
	ERR_CANNOT_CONNECT = utils.Errore{
		ID:      "CANNOT_CONNECT",
		Message: "Cannot connect to voice channel",
	}
)
