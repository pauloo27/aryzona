package command

import (
	"strings"

	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
)

var (
	commandMap            = map[string]*Command{}
	commandLangMap        = map[string]i18n.LanguageName{}
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
	for _, langName := range i18n.LanguagesName {
		if langName == i18n.DefaultLanguageName {
			continue
		}
		cmdLang := i18n.MustGetCommandDefinition(i18n.MustGetLanguage(langName), command.Name)
		if cmdLang == nil {
			logger.Fatalf("Command %s has no translation for language %s", command.Name, langName)
			return // the logger.Fatal will exit the program, but the compiler doesn't know that
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
