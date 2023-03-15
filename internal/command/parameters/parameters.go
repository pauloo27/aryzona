package parameters

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
)

var ParameterText = &command.CommandParameterType{
	BaseType: TypeString,
	Name:     "text",
	Parser: func(index int, args []string) (any, error) {
		return strings.Join(args[index:], " "), nil
	},
}

var ParameterString = &command.CommandParameterType{
	BaseType: TypeString,
	Name:     "string",
	Parser: func(index int, args []string) (any, error) {
		return args[index], nil
	},
}

var ParameterLowerCasedString = &command.CommandParameterType{
	BaseType: TypeString,
	Name:     "lower cased string",
	Parser: func(index int, args []string) (any, error) {
		return strings.ToLower(args[index]), nil
	},
}

var ParameterBool = &command.CommandParameterType{
	BaseType: TypeBool,
	Name:     "bool",
	Parser: func(index int, args []string) (any, error) {
		b, err := strconv.ParseBool(args[index])
		if err != nil {
			return nil, errors.New("invalid boolean value (expected true or false)")
		}
		return b, err
	},
}

var ParameterInt = &command.CommandParameterType{
	BaseType: TypeInt,
	Name:     "int",
	Parser: func(index int, args []string) (any, error) {
		i, err := strconv.Atoi(args[index])
		if err != nil {
			return nil, errors.New("invalid integer number")
		}
		return i, err
	},
}
