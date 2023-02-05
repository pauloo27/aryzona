package tools

import "github.com/Pauloo27/aryzona/internal/command"

var Tools = command.CommandCategory{
	Name:  "tools",
	Emoji: "🔧",
	Commands: []*command.Command{
		&UUIDCommand, &CPFCommand, &CNPJCommand, &PasswordCommand,
	},
}

func init() {
	command.RegisterCategory(Tools)
}
