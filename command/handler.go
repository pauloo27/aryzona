package command

import (
	"github.com/Pauloo27/aryzona/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleCommand(commandName string, args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	command, ok := commandMap[commandName]
	if !ok {
		return
	}

	ctx := &CommandContext{m.Message, m, s, args, nil}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			ctx.Error(utils.Fmt("This command requires `%s`", command.Permission.Name))
			return
		}
	}

	var syntaxError string
	values := []interface{}{}
	if command.Arguments != nil && len(command.Arguments) != 0 {
		parameters := args
		parametersCount := len(parameters)
		for i, argument := range command.Arguments {
			if i >= parametersCount {
				if argument.Required {
					if argument.RequiredMessage == "" {
						syntaxError = utils.Fmt("Argument %s (type %s) missing", argument.Name, argument.Type.Name)
					} else {
						syntaxError = argument.RequiredMessage
					}
				}
			} else {
				// TODO: parse and check
				value, err := argument.Type.Parser(i, args)
				if err != nil {
					syntaxError = err.Error()
				}
				values = append(values, value)
			}
		}
	}
	if syntaxError != "" {
		ctx.Error(syntaxError)
		return
	}
	ctx.Args = values

	command.Handler(ctx)
}
