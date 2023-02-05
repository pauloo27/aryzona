package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
)

var Audio = command.CommandCategory{
	Name:  "audio",
	Emoji: "🎵",
	Commands: []*command.Command{
		&LyricCommand, &RadioCommand, &PlayingCommand, &StopCommand, &PlayCommand,
		&SkipCommand, &PauseCommand, &ResumeCommand, &ShuffleCommand,
	},
}

func init() {
	command.RegisterCategory(Audio)
}
