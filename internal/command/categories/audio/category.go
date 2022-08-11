package audio

import (
	"github.com/Pauloo27/aryzona/internal/command"
)

var Audio = command.CommandCategory{
	Name:  "Audio related stuff",
	Emoji: "🎵",
	Commands: []*command.Command{
		&LyricCommand, &RadioCommand, &PlayingCommand, &StopCommand, &PlayCommand,
		&SkipCommand, &PauseCommand, &ShuffleCommand,
	},
}

func init() {
	command.RegisterCategory(Audio)
}
