package sysmon

import (
	"runtime"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/utils"
)

var Sys = command.Command{
	Name:        "sys",
	Description: "Show system info",
	Handler: func(ctx *command.CommandContext) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		status := utils.Fmt(":computer: %s %s %s",
			runtime.GOOS, runtime.GOARCH, runtime.Version(),
		)
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, status)
	},
}
