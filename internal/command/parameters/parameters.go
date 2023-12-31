package parameters

import (
	"errors"
	"strconv"
	"strings"

	"github.com/pauloo27/aryzona/internal/command"
)

var ParameterText = &command.ParameterType{
	BaseType: TypeString,
	Name:     "text",
	Parser: func(ctx *command.Context, index int, args []string) (any, error) {
		return strings.Join(args[index:], " "), nil
	},
}

var ParameterString = &command.ParameterType{
	BaseType: TypeString,
	Name:     "string",
	Parser: func(ctx *command.Context, index int, args []string) (any, error) {
		return args[index], nil
	},
}

var ParameterLowerCasedString = &command.ParameterType{
	BaseType: TypeString,
	Name:     "lower cased string",
	Parser: func(ctx *command.Context, index int, args []string) (any, error) {
		return strings.ToLower(args[index]), nil
	},
}

var ParameterBool = &command.ParameterType{
	BaseType: TypeBool,
	Name:     "bool",
	Parser: func(ctx *command.Context, index int, args []string) (any, error) {
		b, err := strconv.ParseBool(args[index])
		if err != nil {
			return nil, errors.New(
				ctx.Lang.Validations.ParametersValidations.InvalidBool.Str(),
			)
		}
		return b, err
	},
}

var ParameterInt = &command.ParameterType{
	BaseType: TypeInt,
	Name:     "int",
	Parser: func(ctx *command.Context, index int, args []string) (any, error) {
		i, err := strconv.Atoi(args[index])
		if err != nil {
			return nil, errors.New(
				ctx.Lang.Validations.ParametersValidations.InvalidInt.Str(),
			)
		}
		return i, err
	},
}
