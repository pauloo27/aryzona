package command

import (
	"log/slog"
	"os"
	"strings"

	"github.com/pauloo27/aryzona/internal/i18n"
)

var (
	commandMap            = map[string]*Command{}
	commandLangMap        = map[string]i18n.LanguageName{}
	commandInteractionMap = map[string]*Context{}
	commandList           []*Command
)

var Prefix string

func RegisterCommand(command *Command) {
	commandList = append(commandList, command)
	commandMap[strings.ToLower(command.Name)] = command
	for _, alias := range command.Aliases {
		commandMap[strings.ToLower(alias)] = command
	}
	for _, langName := range i18n.LanguagesName {
		if langName == i18n.DefaultLanguageName {
			continue
		}
		cmdLang := i18n.MustGetCommandDefinition(i18n.MustGetLanguage(langName), command.Name)
		if cmdLang == nil {
			slog.Error("Command translation not found", "commandName", command.Name, "langName", langName)
			os.Exit(-1)
		}
		cmdName := strings.ToLower(cmdLang.Name.Str())
		if _, found := commandMap[cmdName]; found {
			continue
		}
		commandMap[cmdName] = command
		commandLangMap[cmdName] = langName
	}
}

func GetCommandList() []*Command {
	return commandList
}

func GetCommandMap() map[string]*Command {
	return commandMap
}

func RegisterCategory(category Category) {
	if category.OnLoad != nil {
		category.OnLoad()
	}
	if category.Name == "" {
		slog.Error("One category has no name")
		os.Exit(1)
	}
	if category.Emoji == "" {
		slog.Error("Category has no emoji", "categoryName", category.Name)
		os.Exit(1)
	}
	for _, cmd := range category.Commands {
		cmd.category = &category
		RegisterCommand(cmd)
	}
}

func RemoveInteractionHandler(baseID string) {
	delete(commandInteractionMap, baseID)
}
