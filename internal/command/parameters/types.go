package parameters

import "github.com/Pauloo27/aryzona/internal/command"

var (
	TypeString = &command.BaseType{Name: "string"}
	TypeInt    = &command.BaseType{Name: "int"}
	TypeBool   = &command.BaseType{Name: "bool"}
	// TODO: float?
)
