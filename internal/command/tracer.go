package command

import (
	"context"

	"github.com/pauloo27/aryzona/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func startCommandTrace(
	command *Command,
	trigger *TriggerEvent,
	executionID string,
) (context.Context, trace.Span) {
	trCtx, span := tracing.Tracer.Start(context.Background(), "HandleCommand")
	span.SetAttributes(
		attribute.KeyValue{
			Key:   "executionID",
			Value: attribute.StringValue(executionID),
		},
		attribute.KeyValue{
			Key:   "command.name",
			Value: attribute.StringValue(command.Name),
		},
		attribute.KeyValue{
			Key:   "author.id",
			Value: attribute.StringValue(trigger.AuthorID),
		},
		attribute.KeyValue{
			Key:   "guild.id",
			Value: attribute.StringValue(trigger.GuildID),
		},
		attribute.KeyValue{
			Key:   "channel.id",
			Value: attribute.StringValue(trigger.Channel.ID()),
		},
		attribute.KeyValue{
			Key:   "trigger.type",
			Value: attribute.StringValue(string(trigger.Type)),
		},
		attribute.KeyValue{
			Key:   "trigger.message.id",
			Value: attribute.StringValue(trigger.MessageID),
		},
	)

	return trCtx, span
}

func startChildSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracing.Tracer.Start(ctx, name)
}

func addEventToSpan(span trace.Span, name string) {
	span.AddEvent(name)
}
