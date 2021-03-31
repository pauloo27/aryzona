package command

import (
	"github.com/Pauloo27/aryzona/utils"
	"github.com/bwmarrin/discordgo"
)

type CommandContext struct {
	Message       *discordgo.Message
	MessageCreate *discordgo.MessageCreate
	Session       *discordgo.Session
}

type CommandHandler func(*CommandContext)

type Command struct {
	Name, Description string
	Aliases           []string
	Handler           CommandHandler
}

func (ctx *CommandContext) Success(message string) {
	ctx.Session.ChannelMessageSend(
		ctx.Message.ChannelID, utils.Fmt(":green_square: %s", message),
	)
}

func (ctx *CommandContext) Error(message string) {
	ctx.Session.ChannelMessageSend(
		ctx.Message.ChannelID, utils.Fmt(":red_square: %s", message),
	)
}
