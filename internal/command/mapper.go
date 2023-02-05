package command

import (
	"strings"

	"github.com/Pauloo27/logger"
)

var (
	commandMap            = map[string]*Command{}
	commandInteractionMap = map[string]*CommandContext{}
	commandList           []*Command
)

var Prefix string

func RegisterCommand(command *Command) {
	commandList = append(commandList, command)
	commandMap[strings.ToLower(command.Name)] = command
	for _, alias := range command.Aliases {
		commandMap[strings.ToLower(alias)] = command
	}
}

func GetCommandList() []*Command {
	return commandList
}

func GetCommandMap() map[string]*Command {
	return commandMap
}

func RegisterCategory(category CommandCategory) {
	if category.OnLoad != nil {
		category.OnLoad()
	}
	if category.Name == "" {
		logger.Fatal("One category has no name")
	}
	if category.Emoji == "" {
		logger.Fatalf("Category %s has no emoji", category.Name)
	}
	for _, cmd := range category.Commands {
		cmd.category = &category
		RegisterCommand(cmd)
	}
}

func RemoveInteractionHandler(baseID string) {
	delete(commandInteractionMap, baseID)
}
