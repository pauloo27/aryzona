package fun

import "github.com/Pauloo27/aryzona/internal/command"

var Fun = command.CommandCategory{
	Name:  "Fun",
	Emoji: "🎉",
	Commands: []*command.Command{
		&DrawCommand, &EvenCommand, &RollCommand, &ScoreCommand, &XkcdCommand,
		&NewsCommand,
	},
}

func init() {
	command.RegisterCategory(Fun)
}
