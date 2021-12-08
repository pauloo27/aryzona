package discord

type Channel struct {
	ID    string
	Guild *Guild
}

type VoiceChannel struct {
	Channel
}
