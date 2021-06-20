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
	ctx.Session.ChannelMessageSendReply(
		ctx.Message.ChannelID, utils.Fmt(":green_square: %s", message),
		ctx.Message.Reference(),
	)
}

func (ctx *CommandContext) Error(message string) {
	ctx.Session.ChannelMessageSendReply(
		ctx.Message.ChannelID, utils.Fmt(":red_square: %s", message),
		ctx.Message.Reference(),
	)
}

func (ctx *CommandContext) SuccesEmbed(embed *discordgo.MessageEmbed) {
	embed.Color = 0x50fa7b
	ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Reference: ctx.Message.Reference(),
		Embed:     embed,
	})
}

func (ctx *CommandContext) ErrorEmbed(embed *discordgo.MessageEmbed) {
	embed.Color = 0xff5555
	ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Reference: ctx.Message.Reference(),
		Embed:     embed,
	})
}
