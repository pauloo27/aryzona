package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/categories/audio/play"
)

var Audio = command.Category{
	Name:  "audio",
	Emoji: "🎵",
	Commands: []*command.Command{
		&LyricCommand, &RadioCommand, &PlayingCommand, &StopCommand, &play.PlayCommand,
		&SkipCommand, &PauseCommand, &ResumeCommand, &ShuffleCommand, &VolumeCommand,
	},
}

func init() {
	//command.RegisterCategory(Audio)
}
