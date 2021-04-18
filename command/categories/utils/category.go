package utils

import "github.com/Pauloo27/aryzona/command"

var Utils = command.Category{
	Name:     "Utils",
	Commands: []*command.Command{&PingCommand, &UptimeCommand, &HelpCommand},
}
