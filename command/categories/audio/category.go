package audio

import "github.com/Pauloo27/aryzona/command"

var Audio = command.Category{
	Name:     "Audio related stuff",
	Commands: []*command.Command{&RadioCommand, &PlayingCommand},
}
