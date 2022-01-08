package parameters

import (
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
)

var ParameterText = &command.CommandParameterType{
	BaseType: TypeString,
	Name:     "text",
	Parser: func(index int, args []string) (interface{}, error) {
		return strings.Join(args[index:], " "), nil
	},
}

var ParameterString = &command.CommandParameterType{
	BaseType: TypeString,
	Name:     "string",
	Parser: func(index int, args []string) (interface{}, error) {
		return args[index], nil
	},
}

var ParameterBool = &command.CommandParameterType{
	BaseType: TypeBool,
	Name:     "bool",
	Parser: func(index int, args []string) (interface{}, error) {
		return strconv.ParseBool(args[index])
	},
}

var ParameterInt = &command.CommandParameterType{
	BaseType: TypeInt,
	Name:     "int",
	Parser: func(index int, args []string) (interface{}, error) {
		return strconv.Atoi(args[index])
	},
}
