package utils

import "github.com/Pauloo27/aryzona/internal/command"

var Utils = command.CommandCategory{
	Name:  "Utils",
	Emoji: "🔧",
	Commands: []*command.Command{
		&PingCommand, &UptimeCommand, &HelpCommand, &DogCommand, &CatCommand,
		&FoxCommand, &SourceCommand, &XkcdCommand, &ScoreCommand,
		&EvenCommand, &RollCommand, &DrawCommand, &NewsCommand, &UUIDCommand,
		&CPFCommand, &CNPJCommand,
	},
}

func init() {
	command.RegisterCategory(Utils)
}
