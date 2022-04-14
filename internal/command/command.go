package command

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/logger"
)

type CommandHandler func(*CommandContext)
type CommandPermissionChecker func(*CommandContext) bool
type CommandValidationChecker func(*CommandContext) (bool, string)

type CommandContext struct {
	startDate         time.Time
	RawArgs           []string
	Args              []interface{}
	Bot               discord.BotAdapter
	AuthorID, GuildID string
	UsedName          string
	Locals            map[string]interface{}
	Reply             func(string) error
	ReplyEmbed        func(*discord.Embed) error
	Command           *Command
}

type CommandPermission struct {
	Name    string
	Checker CommandPermissionChecker
}

type CommandValidation struct {
	DependsOn   []*CommandValidation
	Description string
	Checker     CommandValidationChecker
}

type CommandParameterTypeParser func(index int, args []string) (interface{}, error)

type BaseType struct {
	Name string
}

type CommandParameterType struct {
	Name     string
	BaseType *BaseType
	Parser   CommandParameterTypeParser
}

type CommandParameter struct {
	ValidValues     []interface{}
	Name            string
	Description     string
	RequiredMessage string
	Type            *CommandParameterType
	ValidValuesFunc func() []interface{}
	Required        bool
}

func (param *CommandParameter) GetValidValues() []interface{} {
	if param.ValidValues != nil {
		return param.ValidValues
	}
	if param.ValidValuesFunc != nil {
		return param.ValidValuesFunc()
	}
	return nil
}

type Command struct {
	Validations       []*CommandValidation
	Parameters        []*CommandParameter
	Aliases           []string
	SubCommands       []*Command
	Name, Description string
	Handler           CommandHandler
	Permission        *CommandPermission
	category          *CommandCategory
	Deferred          bool
}

func (c *Command) GetCategory() *CommandCategory {
	return c.category
}

func (ctx *CommandContext) handleCannotSendMessage(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func (ctx *CommandContext) Success(message string) {
	ctx.handleCannotSendMessage(ctx.SuccessReturning(message))
}

func (ctx *CommandContext) SuccessReturning(message string) error {
	return ctx.Reply(fmt.Sprintf(":green_square: %s", message))
}

func (ctx *CommandContext) Error(message string) {
	ctx.handleCannotSendMessage(ctx.ErrorReturning(message))
}

func (ctx *CommandContext) ErrorReturning(message string) error {
	return ctx.Reply(fmt.Sprintf(":red_square: %s", message))
}

func (ctx *CommandContext) Embed(embed *discord.Embed) {
	ctx.handleCannotSendMessage(ctx.EmbedReturning(embed))
}

func (ctx *CommandContext) EmbedReturning(embed *discord.Embed) error {
	processTime := time.Since(ctx.startDate).Truncate(time.Second)
	duration := fmt.Sprintf("Took %v", processTime)
	if embed.Footer != "" {
		embed.Footer = fmt.Sprintf("%s • %s", embed.Footer, duration)
	} else {
		embed.Footer = duration
	}
	return ctx.ReplyEmbed(embed)
}

func (ctx *CommandContext) SuccessEmbed(embed *discord.Embed) {
	ctx.handleCannotSendMessage(ctx.SuccessEmbedReturning(embed))
}

func (ctx *CommandContext) SuccessEmbedReturning(embed *discord.Embed) error {
	embed.Color = 0x50fa7b
	return ctx.EmbedReturning(embed)
}

func (ctx *CommandContext) ErrorEmbed(embed *discord.Embed) {
	ctx.handleCannotSendMessage(ctx.ErrorEmbedReturning(embed))
}

func (ctx *CommandContext) ErrorEmbedReturning(embed *discord.Embed) error {
	embed.Color = 0xff5555
	return ctx.EmbedReturning(embed)
}
