package command

import "strconv"

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
