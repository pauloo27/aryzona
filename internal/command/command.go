package command

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/logger"
	"github.com/matoous/go-nanoid/v2"
)

type CommandHandler func(*CommandContext)
type CommandPermissionChecker func(*CommandContext) bool
type CommandValidationChecker func(*CommandContext) (bool, string)

type CommandTrigger string

const (
	CommandTriggerSlash   CommandTrigger = "SLASH"
	CommandTriggerMessage CommandTrigger = "MESSAGE"
)

const (
	InteractionBaseIDLength = 10
)

const (
	SuccessEmbedColor = 0x50fa7b
	ErrorEmbedColor   = 0xff5555
)

type CommandContext struct {
	interactionHandler func(id string)
	startDate          time.Time
	RawArgs            []string
	Args               []interface{}
	Bot                discord.BotAdapter
	Member             model.Member
	Channel            model.TextChannel
	AuthorID, GuildID  string
	UsedName           string
	Locals             map[string]interface{}
	Reply              func(string) error
	ReplyEmbed         func(*model.Embed) error
	ReplyComplex       func(*model.ComplexMessage) error
	Edit               func(string) error
	EditEmbed          func(*model.Embed) error
	Command            *Command
	Trigger            CommandTrigger
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

func (ctx *CommandContext) Successf(format string, a ...any) {
	ctx.Success(fmt.Sprintf(format, a...))
}

func (ctx *CommandContext) SuccessReturning(message string) error {
	return ctx.SuccessEmbedReturning(model.NewEmbed().WithDescription(message))
}

func (ctx *CommandContext) Error(message string) {
	ctx.handleCannotSendMessage(ctx.ErrorReturning(message))
}

func (ctx *CommandContext) Errorf(format string, a ...any) {
	ctx.Error(fmt.Sprintf(format, a...))
}

func (ctx *CommandContext) ErrorReturning(message string) error {
	return ctx.ErrorEmbedReturning(model.NewEmbed().WithDescription(message))
}

func (ctx *CommandContext) Embed(embed *model.Embed) {
	ctx.handleCannotSendMessage(ctx.EmbedReturning(embed))
}

func (ctx *CommandContext) EmbedReturning(embed *model.Embed) error {
	ctx.AddCommandDuration(embed)
	return ctx.ReplyEmbed(embed)
}

func (ctx *CommandContext) SuccessEmbed(embed *model.Embed) {
	ctx.handleCannotSendMessage(ctx.SuccessEmbedReturning(embed))
}

func (ctx *CommandContext) SuccessEmbedReturning(embed *model.Embed) error {
	embed.Color = SuccessEmbedColor
	return ctx.EmbedReturning(embed)
}

func (ctx *CommandContext) ErrorEmbed(embed *model.Embed) {
	ctx.handleCannotSendMessage(ctx.ErrorEmbedReturning(embed))
}

func (ctx *CommandContext) ErrorEmbedReturning(embed *model.Embed) error {
	embed.Color = ErrorEmbedColor
	return ctx.EmbedReturning(embed)
}

func (ctx *CommandContext) AddCommandDuration(embed *model.Embed) {
	processTime := time.Since(ctx.startDate).Truncate(time.Second)
	duration := fmt.Sprintf("Took %v", processTime)
	if embed.Footer != "" {
		embed.Footer = fmt.Sprintf("%s • %s", embed.Footer, duration)
	} else {
		embed.Footer = duration
	}
}

func (ctx *CommandContext) RegisterInteractionHandler(interactionHandler func(id string)) (baseID string, err error) {
	for {
		baseID, err = gonanoid.New(InteractionBaseIDLength)
		if err != nil {
			logger.Error(err)
			return "", err
		}
		if _, found := commandInteractionMap[baseID]; !found {
			break
		}
	}
	ctx.interactionHandler = interactionHandler
	commandInteractionMap[baseID] = ctx
	return baseID, nil
}
