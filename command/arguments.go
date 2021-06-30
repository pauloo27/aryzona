package command

import (
	"strconv"

	"github.com/Pauloo27/aryzona/utils"
)

var ArgumentString = &CommandArgumentType{
	Name: "string",
	Parser: func(index int, args []string) (interface{}, error) {
		return args[index], nil
	},
}

var ArgumentInt = &CommandArgumentType{
	Name: "int",
	Parser: func(index int, args []string) (interface{}, error) {
		return strconv.Atoi(args[index])
	},
}

func ErrRequiredArgument(argument *CommandArgument) *utils.Errore {
	var message string
	if argument != nil {
		if argument.RequiredMessage != "" {
			message = argument.RequiredMessage
		} else {
			message = utils.Fmt("Argument %s (type %s) missing", argument.Name, argument.Type.Name)
		}
	}
	return &utils.Errore{
		ID:      "REQUIRED_ARGUMENT_MISSING",
		Message: message,
	}
}

func ErrInvalidValue(argument *CommandArgument, validValues []interface{}) *utils.Errore {
	var message string
	if argument != nil {
		message = utils.Fmt("Invalid value for %s", argument.Name)
		if validValues != nil {
			message += utils.Fmt(". Valid values are %v", validValues)
		}
	}
	return &utils.Errore{
		ID:      "INVALID_VALUE",
		Message: message,
	}
}

func ErrCannotParseArgument(argument *CommandArgument, err error) *utils.Errore {
	var message string
	if err != nil {
		message = err.Error()
	}
	return &utils.Errore{
		ID:      "CANNOT_PARSE_ARGUMENT",
		Message: message,
	}
}

func (command *Command) ValidateArguments(args []string) (values []interface{}, syntaxError *utils.Errore) {
	if command.Arguments != nil && len(command.Arguments) != 0 {
		parameters := args
		parametersCount := len(parameters)
		for i, argument := range command.Arguments {
			if i >= parametersCount {
				if argument.Required {
					syntaxError = ErrRequiredArgument(argument)
					break
				}
			} else {
				value, err := argument.Type.Parser(i, args)
				if err != nil {
					syntaxError = ErrCannotParseArgument(argument, err)
					break
				}
				if command.ValidValues != nil {
					valid := false
					for _, v := range command.ValidValues {
						if v == value {
							valid = true
							break
						}
					}
					if !valid {
						syntaxError = ErrInvalidValue(argument, command.ValidValues)
						break
					}
				}
				values = append(values, value)
			}
		}
	}
	return
}
