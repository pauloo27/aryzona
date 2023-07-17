package command

import (
	"github.com/pauloo27/aryzona/internal/discord/model"
)

type Result struct {
	Success bool
	Message *model.ComplexMessage
}
type Handler func(*Context) Result
type PermissionChecker func(*Context) bool
type ValidationChecker func(*Context) (bool, string)

type TriggerType string

const (
	CommandTriggerSlash   TriggerType = "SLASH"
	CommandTriggerMessage TriggerType = "MESSAGE"
)

const (
	InteractionBaseIDLength = 10
)

type InteractionHandler func(id, userID, baseID string) (newMessage *model.ComplexMessage, done bool)

type Permission struct {
	Name    string
	Checker PermissionChecker
}

type Validation struct {
	Name      string
	DependsOn []*Validation
	Checker   ValidationChecker
}

type ParameterTypeParser func(ctx *Context, index int, args []string) (any, error)

type BaseType struct {
	Name string
}

type ParameterType struct {
	Name     string
	BaseType *BaseType
	Parser   ParameterTypeParser
}

type Parameter struct {
	ValidValues     []any
	Name            string
	Type            *ParameterType
	ValidValuesFunc func() []any
	Required        bool
}

func (param *Parameter) GetValidValues() []any {
	if param.ValidValues != nil {
		return param.ValidValues
	}
	if param.ValidValuesFunc != nil {
		return param.ValidValuesFunc()
	}
	return nil
}

type Command struct {
	parent *Command

	Validations []*Validation
	Parameters  []*Parameter
	Aliases     []string
	SubCommands []*Command
	Name        string
	Handler     Handler
	Permission  *Permission
	category    *Category
	Deferred    bool
	Ephemeral   bool
}

func (c *Command) GetCategory() *Category {
	return c.category
}
