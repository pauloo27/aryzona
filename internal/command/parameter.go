package command

import (
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/aryzona/internal/utils/errore"
)

func ErrRequiredParameter(param *CommandParameter) *errore.Errore {
	var message string
	if param != nil {
		if param.RequiredMessage != "" {
			message = param.RequiredMessage
		} else {
			message = utils.Fmt("Parameter %s (type %s) missing", param.Description, param.Type.Name)
		}
	}
	return &errore.Errore{
		ID:      "REQUIRED_PARAMETER_MISSING",
		Message: message,
	}
}

func ErrInvalidValue(param *CommandParameter) *errore.Errore {
	var message string
	if param != nil {
		message = utils.Fmt("Invalid value for %s. Valid  values are: %s", param.Description, param.GetValidValues())
	}
	return &errore.Errore{
		ID:      "INVALID_VALUE",
		Message: message,
	}
}

func ErrCannotParseParameter(argument *CommandParameter, err error) *errore.Errore {
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
 ARGUMENTS -> what was passed to the command
 PARAMETERS -> what the command expect

 this function is supposed to, given a string array (the arguments),
 tell if a command can be executed based on it's parameter definition. The
 function checks for "required" parameters, acceptable values (when only a
 pre-defined set of values are accepted) and also for type validation
 (eg, if a argument is an integer number but the parameter is a invalid string).

 If no errors are returned, then we are good to go.
*/
func (command *Command) ValidateParameters(parameters []string) (values []interface{}, syntaxError *errore.Errore) {

	if command.Parameters == nil || len(command.Parameters) == 0 {
		return
	}

	parametersCount := len(parameters)

	for i, parameter := range command.Parameters {
		if i >= parametersCount {
			if parameter.Required {
				syntaxError = ErrRequiredParameter(parameter)
			}
			break
		}

		value, err := parameter.Type.Parser(i, parameters)
		if err != nil {
			syntaxError = ErrCannotParseParameter(parameter, err)
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
				syntaxError = ErrInvalidValue(parameter)
				break
			}
		}

		values = append(values, value)
	}

	return
}
