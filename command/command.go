package command

import (
	"github.com/Pauloo27/aryzona/utils"
	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(*CommandContext)
type CommandPermissionChecker func(*CommandContext) bool

type CommandContext struct {
	Message       *discordgo.Message
	MessageCreate *discordgo.MessageCreate
	Session       *discordgo.Session
	Args          []string
}

type CommandPermission struct {
	Name    string
	Checker CommandPermissionChecker
}

type Command struct {
	Name, Description string
	Aliases           []string
	Handler           CommandHandler
	Permission        *CommandPermission
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
