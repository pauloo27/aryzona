package command

import (
	"errors"
	"fmt"

	"github.com/Pauloo27/aryzona/internal/utils"
)

var (
	ErrRequireParameter    = errors.New("required parameter missing")
	ErrCannotParseArgument = errors.New("cannot parse argument")
	ErrInvalidValue        = errors.New("invalid value")
)

func NewErrRequiredParameter(param *CommandParameter) error {
	var message string
	if param != nil {
		if param.RequiredMessage != "" {
			message = param.RequiredMessage
		} else {
			message = utils.Fmt("parameter `%s` (type %s) missing", param.Description, param.Type.Name)
		}
	}
	return fmt.Errorf("%w: %s", ErrRequireParameter, message)
}

func NewErrInvalidValue(param *CommandParameter) error {
	var message string
	if param != nil {
		message = utils.Fmt("invalid value for `%s`. Valid  values are: `%s`", param.Description, param.GetValidValues())
	}
	return fmt.Errorf("%w: %s", ErrInvalidValue, message)
}

func NewErrCannotParseParameter(argument *CommandParameter, err error) error {
	var message string
	if err != nil {
		message = err.Error()
	}
	return fmt.Errorf("%w: %s", ErrCannotParseArgument, message)
}

/*
 ARGUMENTS -> what was passed to the command
 PARAMETERS -> what the command expect

 this function is supposed to, given a string array (the arguments),
 tell if a command can be executed based on it's parameter definition. The
 function checks for "required" parameters, acceptable values (when only a
 pre-defined set of values are accepted) and also for type validation
 (eg, if a argument is an integer number but the parameter is a invalid string).

 If no errors are returned, then we are good to go.
*/
func (command *Command) ValidateParameters(parameters []string) (values []interface{}, syntaxError error) {

	if command.Parameters == nil || len(command.Parameters) == 0 {
		return
	}

	parametersCount := len(parameters)

	for i, parameter := range command.Parameters {
		if i >= parametersCount {
			if parameter.Required {
				syntaxError = NewErrRequiredParameter(parameter)
			}
			break
		}

		value, err := parameter.Type.Parser(i, parameters)
		if err != nil {
			syntaxError = NewErrCannotParseParameter(parameter, err)
			break
		}

		if parameter.GetValidValues() != nil {
			valid := false
			for _, v := range parameter.GetValidValues() {
				if v == value {
					valid = true
					break
				}
			}
			if !valid {
				syntaxError = NewErrInvalidValue(parameter)
				break
			}
		}

		values = append(values, value)
	}

	return
}
