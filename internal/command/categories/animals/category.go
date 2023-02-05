package animals

import "github.com/Pauloo27/aryzona/internal/command"

var Animals = command.CommandCategory{
	Name:  "animals",
	Emoji: "🐕",
	Commands: []*command.Command{
		&DogCommand, &CatCommand, &FoxCommand,
	},
}

func init() {
	command.RegisterCategory(Animals)
}
