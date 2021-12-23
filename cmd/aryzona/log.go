package main

import (
	"fmt"
	"runtime/debug"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/logger"
)

func init() {
	logger.AddLogListener(func(level logger.Level, params ...interface{}) {
		if discord.Bot.StartedAt() == nil || (level != logger.ERROR && level != logger.FATAL) {
			return
		}

		c, err := discord.OpenChatWithOwner()
		if err != nil {
			// to avoid loops, do not call the logger again
			fmt.Println("Cannot open chat with owner", err)
			return
		}

		embed := discord.NewEmbed().
			WithFieldInline("Message", fmt.Sprintln(params...)).
			WithDescription(utils.Fmt("```go\n%s\n```", string(debug.Stack()))).
			WithColor(0xff5555).
			WithTitle(utils.Fmt("Oops! [%s]", level.Name))

		_, err = discord.Bot.SendEmbedMessage(c.ID(), embed)
		if err != nil {
			fmt.Println("Cannot log to Discord", err)
			return
		}
	})
}
