package voicer

type VoicerError struct {
	ID, Message string
}

func (v VoicerError) Error() string {
	return v.Message
}

var (
	ERR_ALREADY_PLAYING = VoicerError{"ALREADY_PLAYING", "Already playing something in the current voice channel"}
)

func IsVoicerError(err error) (bool, VoicerError) {
	e, ok := err.(VoicerError)
	return ok, e
}

func Is(err error, vErr VoicerError) bool {
	voicerError, ok := err.(VoicerError)
	if !ok {
		return false
	}

	return voicerError.ID == vErr.ID
}
