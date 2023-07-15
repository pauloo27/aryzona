package command

import (
	"errors"
	"fmt"

	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
)

var (
	ErrRequireParameter    = errors.New("required parameter missing")
	ErrCannotParseArgument = errors.New("cannot parse argument")
	ErrInvalidValue        = errors.New("invalid value")
)

func (command *Command) ValidateParameters(ctx *Context) (values []any, syntaxError error) {
	arguments := ctx.RawArgs
	lang := ctx.Lang

	if command.Parameters == nil || len(command.Parameters) == 0 {
		return
	}

	argumentsCount := len(arguments)

	for i, parameter := range command.Parameters {
		if i >= argumentsCount {
			if parameter.Required {
				syntaxError = NewErrRequiredParameter(lang, command, parameter)
			}
			break
		}

		value, err := parameter.Type.Parser(ctx, i, arguments)
		if err != nil {
			syntaxError = NewErrCannotParseParameter(lang, err)
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
				syntaxError = NewErrInvalidValue(lang, command, parameter)
				break
			}
		}

		values = append(values, value)
	}

	return
}

func NewErrRequiredParameter(lang *i18n.Language, cmd *Command, param *Parameter) error {
	localizedParamName := mustGetLocalizedParamName(lang, cmd, param)

	var message string
	if param != nil {
		message = lang.Validations.ParametersValidations.RequiredParam.Str(localizedParamName)
	}
	return fmt.Errorf("%w: %s", ErrRequireParameter, message)
}

func NewErrInvalidValue(lang *i18n.Language, cmd *Command, param *Parameter) error {
	localizedParamName := mustGetLocalizedParamName(lang, cmd, param)

	var message string
	if param != nil {
		message = lang.Validations.ParametersValidations.InvalidValue.Str(localizedParamName, param.GetValidValues())
	}
	return fmt.Errorf("%w: %s", ErrInvalidValue, message)
}

func NewErrCannotParseParameter(lang *i18n.Language, err error) error {
	var message string
	if err != nil {
		message = err.Error()
	}
	return fmt.Errorf("%w: %s", ErrCannotParseArgument, message)
}

func mustGetLocalizedParamName(lang *i18n.Language, cmd *Command, param *Parameter) string {
	var paramIdx int
	for i, p := range cmd.Parameters {
		if p.Name == param.Name {
			paramIdx = i
		}
	}

	var localizedParamName any
	var err error

	if cmd.parent == nil {
		localizedParamName, err = lang.RawMap.Get("commands", cmd.Name, "definition", "parameters", paramIdx, "name")
	} else {

		var subCommandIdx int

		for i, subCmd := range cmd.parent.SubCommands {
			if subCmd.Name == cmd.Name {
				subCommandIdx = i
			}
		}

		localizedParamName, err = lang.RawMap.Get("commands", cmd.parent.Name, "definition", "subCommands", subCommandIdx, "parameters", paramIdx, "name")
	}

	if err != nil {
		logger.Errorf("Failed to get localized parameter name: %v", err)
		return param.Name
	}

	return localizedParamName.(string)
}
