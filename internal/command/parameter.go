package command

import (
	"errors"
	"fmt"
)

var (
	ErrRequireParameter    = errors.New("required parameter missing")
	ErrCannotParseArgument = errors.New("cannot parse argument")
	ErrInvalidValue        = errors.New("invalid value")
)

func (command *Command) ValidateParameters(parameters []string) (values []any, syntaxError error) {
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
			syntaxError = NewErrCannotParseParameter(err)
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

func NewErrRequiredParameter(param *CommandParameter) error {
	var message string
	if param != nil {
		validValues := param.GetValidValues()
		if param.RequiredMessage != "" {
			message = param.RequiredMessage
		} else if validValues != nil {
			message = fmt.Sprintf(
				"parameter `%s` (type %s, valid values are: `%v`) missing",
				param.Name,
				param.Type.Name,
				validValues,
			)
		} else {
			message = fmt.Sprintf("parameter `%s` (type %s) missing", param.Name, param.Type.Name)
		}
	}
	return fmt.Errorf("%w: %s", ErrRequireParameter, message)
}

func NewErrInvalidValue(param *CommandParameter) error {
	var message string
	if param != nil {
		message = fmt.Sprintf("invalid value for `%s`. Valid  values are: `%v`", param.Name, param.GetValidValues())
	}
	return fmt.Errorf("%w: %s", ErrInvalidValue, message)
}

func NewErrCannotParseParameter(err error) error {
	var message string
	if err != nil {
		message = err.Error()
	}
	return fmt.Errorf("%w: %s", ErrCannotParseArgument, message)
}
