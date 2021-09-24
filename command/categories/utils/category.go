package utils

import "github.com/Pauloo27/aryzona/command"

var Utils = command.CommandCategory{
	Name:  "Utils",
	Emoji: "🔧",
	Commands: []*command.Command{
		&PingCommand, &UptimeCommand, &HelpCommand, &WoofCommand, &MeowCommand,
		&FloofCommand, &SourceCommand, &XkcdCommand, &ScoreCommand, &LiveCommand,
		&EvenCommand, &DiceCommand,
	},
}

func init() {
	command.RegisterCategory(Utils)
}
