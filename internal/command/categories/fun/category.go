package fun

import "github.com/pauloo27/aryzona/internal/command"

var Fun = command.CommandCategory{
	Name:  "fun",
	Emoji: "🎉",
	Commands: []*command.Command{
		&PickCommand, &EvenCommand, &RollCommand, &ScoreCommand,
		&NewsCommand, &JokeCommand, &FollowCommand, &UnFollowCommand, &XkcdCommand,
	},
}

func init() {
	command.RegisterCategory(Fun)
}
