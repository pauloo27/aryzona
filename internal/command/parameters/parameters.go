package parameters

import (
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
)

var ParameterText = &command.CommandParameterType{
	Name: "string",
	Parser: func(index int, args []string) (interface{}, error) {
		return strings.Join(args[index:], " "), nil
	},
}

var ParameterString = &command.CommandParameterType{
	Name: "string",
	Parser: func(index int, args []string) (interface{}, error) {
		return args[index], nil
	},
}

var ParameterBool = &command.CommandParameterType{
	Name: "bool",
	Parser: func(index int, args []string) (interface{}, error) {
		return strconv.ParseBool(args[index])
	},
}

var ParameterInt = &command.CommandParameterType{
	Name: "int",
	Parser: func(index int, args []string) (interface{}, error) {
		return strconv.Atoi(args[index])
	},
}
