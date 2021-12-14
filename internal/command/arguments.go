package command

import (
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/aryzona/internal/utils/errore"
)

var ArgumentText = &CommandArgumentType{
	Name: "string",
	Parser: func(index int, args []string) (interface{}, error) {
		return strings.Join(args[index:], " "), nil
	},
}

var ArgumentString = &CommandArgumentType{
	Name: "string",
	Parser: func(index int, args []string) (interface{}, error) {
		return args[index], nil
	},
}

var ArgumentBool = &CommandArgumentType{
	Name: "bool",
	Parser: func(index int, args []string) (interface{}, error) {
		return strconv.ParseBool(args[index])
	},
}

var ArgumentInt = &CommandArgumentType{
	Name: "int",
	Parser: func(index int, args []string) (interface{}, error) {
		return strconv.Atoi(args[index])
	},
}

func ErrRequiredArgument(argument *CommandArgument) *errore.Errore {
	var message string
	if argument != nil {
		if argument.RequiredMessage != "" {
			message = argument.RequiredMessage
		} else {
			message = utils.Fmt("Argument %s (type %s) missing", argument.Description, argument.Type.Name)
		}
	}
	return &errore.Errore{
		ID:      "REQUIRED_ARGUMENT_MISSING",
		Message: message,
	}
}

func ErrInvalidValue(argument *CommandArgument) *errore.Errore {
	var message string
	if argument != nil {
		message = utils.Fmt("Invalid value for %s. Valid  values are: %s", argument.Description, argument.GetValidValues())
	}
	return &errore.Errore{
		ID:      "INVALID_VALUE",
		Message: message,
	}
}

func ErrCannotParseArgument(argument *CommandArgument, err error) *errore.Errore {
	var message string
	if err != nil {
		message = err.Error()
	}
	return &errore.Errore{
		ID:      "CANNOT_PARSE_ARGUMENT",
		Message: message,
	}
}

/*
 PARAMETERS -> what was passed to the command
 ARGUMENTS -> what the command expect

 this function is supposed to, given a string array (the parameters),
 tell if a command can be executed based on it's arguments. The function checks
 for "required" arguments, acceptable values (when only a pre-defined set of
 values are accepted) and also for type validation (eg, if a argument
 is an integer number but the parameter is a invalid string).

 If not errors are returned, then we are good to go.
*/
func (command *Command) ValidateArguments(parameters []string) (values []interface{}, syntaxError *errore.Errore) {

	if command.Arguments == nil || len(command.Arguments) == 0 {
		return
	}

	parametersCount := len(parameters)

	for i, argument := range command.Arguments {
		if i >= parametersCount {
			if argument.Required {
				syntaxError = ErrRequiredArgument(argument)
			}
			break
		}

		value, err := argument.Type.Parser(i, parameters)
		if err != nil {
			syntaxError = ErrCannotParseArgument(argument, err)
			break
		}

		if argument.GetValidValues() != nil {
			valid := false
			for _, v := range argument.GetValidValues() {
				if v == value {
					valid = true
					break
				}
			}
			if !valid {
				syntaxError = ErrInvalidValue(argument)
				break
			}
		}

		values = append(values, value)
	}

	return
}
