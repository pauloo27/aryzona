package command

import (
	"fmt"
	"time"

	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
	gonanoid "github.com/matoous/go-nanoid/v2"
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
	PendingEmbedColor = 0x00add8
)

type InteractionHandler func(fullID, baseID, userID string) (newMessage *model.ComplexMessage, done bool)

type CommandContext struct {
	interactionHandler InteractionHandler
	Lang               *i18n.Language
	T                  any
	startTime          time.Time
	RawArgs            []string
	Args               []any
	Bot                discord.BotAdapter
	Channel            model.TextChannel
	AuthorID, GuildID  string
	UsedName           string
	Locals             map[string]any
	Reply              func(string) error
	ReplyEmbed         func(*model.Embed) error
	ReplyComplex       func(*model.ComplexMessage) error
	EditComplex        func(*model.ComplexMessage) error
	Edit               func(string) error
	EditEmbed          func(*model.Embed) error
	Command            *Command
	Trigger            CommandTrigger

	executionID string
	processTime time.Duration
}

type CommandPermission struct {
	Name    string
	Checker CommandPermissionChecker
}

type CommandValidation struct {
	Name      string
	DependsOn []*CommandValidation
	Checker   CommandValidationChecker
}

type CommandParameterTypeParser func(ctx *CommandContext, index int, args []string) (any, error)

type BaseType struct {
	Name string
}

type CommandParameterType struct {
	Name     string
	BaseType *BaseType
	Parser   CommandParameterTypeParser
}

type CommandParameter struct {
	ValidValues     []any
	Name            string
	Type            *CommandParameterType
	ValidValuesFunc func() []any
	Required        bool
}

func (param *CommandParameter) GetValidValues() []any {
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

	Validations []*CommandValidation
	Parameters  []*CommandParameter
	Aliases     []string
	SubCommands []*Command
	Name        string
	Handler     CommandHandler
	Permission  *CommandPermission
	category    *CommandCategory
	Deferred    bool
	Ephemeral   bool
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
	ctx.processTime = time.Since(ctx.startTime)
	duration := ctx.Lang.Took.Str(ctx.processTime.Truncate(time.Second))
	if embed.Footer != "" {
		embed.Footer = fmt.Sprintf("%s • %s", embed.Footer, duration)
	} else {
		embed.Footer = duration
	}
}

func (ctx *CommandContext) RegisterInteractionHandler(handler InteractionHandler) (baseID string, err error) {
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
	ctx.interactionHandler = handler
	commandInteractionMap[baseID] = ctx
	go func() {
		time.Sleep(5 * time.Minute)
		RemoveInteractionHandler(baseID)
	}()
	return baseID, nil
}
