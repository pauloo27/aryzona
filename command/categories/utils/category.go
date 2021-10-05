package utils

import "github.com/Pauloo27/aryzona/command"

var Utils = command.CommandCategory{
	Name:  "Utils",
	Emoji: "🔧",
	Commands: []*command.Command{
		&PingCommand, &UptimeCommand, &HelpCommand, &DogCommand, &CatCommand,
		&FoxCommand, &SourceCommand, &XkcdCommand, &ScoreCommand, &LiveCommand,
		&EvenCommand, &RollCommand, &UpdateCommand,
	},
}

func init() {
	command.RegisterCategory(Utils)
}
