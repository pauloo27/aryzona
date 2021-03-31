package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/utils"
)

var PingCommand = command.Command{
	Name: "ping", Description: "Get the bot latency",
	Aliases: []string{"p", "pong"},
	Handler: func(ctx *command.CommandContext) {
		ctx.Success(utils.Fmt("Bot latency with Discord servers is %d ms", ctx.Session.HeartbeatLatency().Milliseconds()))
	},
}
