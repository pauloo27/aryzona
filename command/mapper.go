package command

import "strings"

var commandMap = map[string]*Command{}

var Prefix string

func RegisterCommand(command *Command) {
	commandMap[strings.ToLower(command.Name)] = command
	for _, alias := range command.Aliases {
		commandMap[strings.ToLower(alias)] = command
	}
}

func GetCommandMap() map[string]*Command {
	return commandMap
}

func RegisterCategory(category CommandCategory) {
	if category.OnLoad != nil {
		category.OnLoad()
	}
	for _, cmd := range category.Commands {
		cmd.category = &category
		RegisterCommand(cmd)
	}
}
