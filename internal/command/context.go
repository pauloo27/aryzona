package command

import (
	"context"
	"fmt"
	"time"

	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"go.opentelemetry.io/otel/trace"
)

const (
	SuccessEmbedColor = 0x50fa7b
	ErrorEmbedColor   = 0xff5555
	PendingEmbedColor = 0x00add8
)

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
	trCtx       context.Context
	span        trace.Span
	trigger     *TriggerEvent
}

func (ctx *Context) RegisterInteractionHandler(baseID string, handler InteractionHandler) {
	ctx.interactionHandler = handler
	commandInteractionMap[baseID] = ctx
}

func (ctx *Context) Empty() Result {
	return Result{
		Success: true,
		Message: nil,
	}
}

func (ctx *Context) EditComplexMessage(message *model.ComplexMessage) error {
	return ctx.trigger.Edit(ctx, message)
}

func (ctx *Context) appendExecutionInfoToEmbed(embed *model.Embed) []*model.Embed {
	took := ctx.Lang.Took.Str(time.Since(ctx.startTime).Truncate(time.Second))
	extraFooter := fmt.Sprintf("%s • %s", took, ctx.executionID)

	if embed.Footer == "" {
		embed.Footer = extraFooter
	} else {
		embed.Footer = fmt.Sprintf("%s • %s", embed.Footer, extraFooter)
	}
	return []*model.Embed{embed}
}

func (ctx *Context) Embed(embed *model.Embed) Result {
	return Result{
		Success: true,
		Message: &model.ComplexMessage{
			Embeds: ctx.appendExecutionInfoToEmbed(embed),
		},
	}
}

func (ctx *Context) SuccessEmbed(embed *model.Embed) Result {
	embed.Color = SuccessEmbedColor
	return Result{
		Success: true,
		Message: &model.ComplexMessage{
			Embeds: ctx.appendExecutionInfoToEmbed(embed),
		},
	}
}

func (ctx *Context) Success(message string) Result {
	return ctx.SuccessEmbed(&model.Embed{
		Color:       SuccessEmbedColor,
		Description: message,
	})
}

func (ctx *Context) Successf(format string, args ...any) Result {
	return ctx.Success(fmt.Sprintf(format, args...))
}

func (ctx *Context) ErrorEmbed(embed *model.Embed) Result {
	embed.Color = ErrorEmbedColor
	return Result{
		Success: false,
		Message: &model.ComplexMessage{
			Embeds: ctx.appendExecutionInfoToEmbed(embed),
		},
	}
}

func (ctx *Context) Error(message string) Result {
	return ctx.ErrorEmbed(&model.Embed{
		Description: message,
		Color:       ErrorEmbedColor,
	})
}

func (ctx *Context) Errorf(format string, args ...any) Result {
	return ctx.Error(fmt.Sprintf(format, args...))
}

func (ctx *Context) ReplyRaw(content string) Result {
	return Result{
		Success: true,
		Message: &model.ComplexMessage{
			Content: content,
		},
	}
}

func (ctx *Context) ReplyWithInteraction(
	baseID string,
	message *model.ComplexMessage,
	handler InteractionHandler,
) Result {
	ctx.RegisterInteractionHandler(baseID, handler)
	return Result{
		Success: true,
		Message: message,
	}
}
