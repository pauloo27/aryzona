package command

var commandMap = map[string]*Command{}

var Prefix string

func RegisterCommand(command *Command) {
	commandMap[command.Name] = command
	for _, alias := range command.Aliases {
		commandMap[alias] = command
	}
}
