package command

import (
	"fmt"
	"time"

	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
)

type Handler func(*Context)
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

const (
	SuccessEmbedColor = 0x50fa7b
	ErrorEmbedColor   = 0xff5555
	PendingEmbedColor = 0x00add8
)

type InteractionHandler func(id, userID, baseID string) (newMessage *model.ComplexMessage, done bool)

type Context struct {
	interactionHandler           InteractionHandler
	Lang                         *i18n.Language
	T                            any
	startTime                    time.Time
	RawArgs                      []string
	Args                         []any
	Bot                          discord.BotAdapter
	Channel                      model.TextChannel
	MessageID, AuthorID, GuildID string
	UsedName                     string
	Locals                       map[string]any
	Command                      *Command
	TriggerType                  TriggerType

	executionID string
	processTime time.Duration
	trigger     *TriggerEvent
}

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

func (ctx *Context) handleCannotSendMessage(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func (ctx *Context) logResponse() {
	logger.Logf(
		CommandLogLevel,
		"[i %s] got response %s, took %s",
		ctx.executionID,
		"<omitted>", // dont really care about the response
		ctx.processTime,
	)
}

func (ctx *Context) Reply(message string) error {
	ctx.logResponse()

	return ctx.trigger.Reply(ctx, &model.ComplexMessage{
		Content: message,
	})
}

func (ctx *Context) ReplyEmbed(embed *model.Embed) error {
	ctx.logResponse()
	return ctx.trigger.Reply(ctx, &model.ComplexMessage{
		Embeds: []*model.Embed{embed},
	})
}

func (ctx *Context) ReplyComplex(message *model.ComplexMessage) error {
	ctx.logResponse()
	return ctx.trigger.Reply(ctx, message)
}

func (ctx *Context) EditComplex(message *model.ComplexMessage) error {
	ctx.logResponse()
	return ctx.trigger.Edit(ctx, message)
}

func (ctx *Context) Edit(message string) error {
	ctx.logResponse()
	return ctx.trigger.Edit(ctx, &model.ComplexMessage{
		Content: message,
	})
}

func (ctx *Context) EditEmbed(embed *model.Embed) error {
	ctx.logResponse()
	return ctx.trigger.Edit(ctx, &model.ComplexMessage{
		Embeds: []*model.Embed{embed},
	})
}

func (ctx *Context) Success(message string) {
	ctx.handleCannotSendMessage(ctx.SuccessReturning(message))
}

func (ctx *Context) ReplyWithInteraction(
	baseID string,
	message *model.ComplexMessage, handler InteractionHandler,
) error {
	ctx.RegisterInteractionHandler(baseID, handler)
	return ctx.ReplyComplex(message)
}

func (ctx *Context) Successf(format string, a ...any) {
	ctx.Success(fmt.Sprintf(format, a...))
}

func (ctx *Context) SuccessReturning(message string) error {
	return ctx.SuccessEmbedReturning(model.NewEmbed().WithDescription(message))
}

func (ctx *Context) Error(message string) {
	ctx.handleCannotSendMessage(ctx.ErrorReturning(message))
}

func (ctx *Context) Errorf(format string, a ...any) {
	ctx.Error(fmt.Sprintf(format, a...))
}

func (ctx *Context) ErrorReturning(message string) error {
	return ctx.ErrorEmbedReturning(model.NewEmbed().WithDescription(message))
}

func (ctx *Context) Embed(embed *model.Embed) {
	ctx.handleCannotSendMessage(ctx.EmbedReturning(embed))
}

func (ctx *Context) EmbedReturning(embed *model.Embed) error {
	ctx.AddCommandDuration(embed)
	return ctx.ReplyEmbed(embed)
}

func (ctx *Context) SuccessEmbed(embed *model.Embed) {
	ctx.handleCannotSendMessage(ctx.SuccessEmbedReturning(embed))
}

func (ctx *Context) SuccessEmbedReturning(embed *model.Embed) error {
	embed.Color = SuccessEmbedColor
	return ctx.EmbedReturning(embed)
}

func (ctx *Context) ErrorEmbed(embed *model.Embed) {
	ctx.handleCannotSendMessage(ctx.ErrorEmbedReturning(embed))
}

func (ctx *Context) ErrorEmbedReturning(embed *model.Embed) error {
	embed.Color = ErrorEmbedColor
	return ctx.EmbedReturning(embed)
}

func (ctx *Context) AddCommandDuration(embed *model.Embed) {
	ctx.processTime = time.Since(ctx.startTime)
	duration := ctx.Lang.Took.Str(ctx.processTime.Truncate(time.Second))
	if embed.Footer != "" {
		embed.Footer = fmt.Sprintf("%s • %s", embed.Footer, duration)
	} else {
		embed.Footer = duration
	}
}

func (ctx *Context) RegisterInteractionHandler(baseID string, handler InteractionHandler) {
	ctx.interactionHandler = handler
	commandInteractionMap[baseID] = ctx
}
