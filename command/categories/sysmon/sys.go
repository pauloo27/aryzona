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
		status := utils.Fmt(":computer: %s %s %s\n\n:abacus: RAM (MB): %d/%d\n\n",
			runtime.GOOS, runtime.GOARCH, runtime.Version(),
			memStats.Alloc/1024, memStats.Sys/1024,
		)
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, status)
	},
}
