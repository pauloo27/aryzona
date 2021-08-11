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
	RawArgs       []string
	Args          []interface{}
}

type CommandPermission struct {
	Name    string
	Checker CommandPermissionChecker
}

type CommandArgumentTypeParser func(index int, args []string) (interface{}, error)

type CommandArgumentType struct {
	Name   string
	Parser CommandArgumentTypeParser
}

type CommandArgument struct {
	Name            string
	Type            *CommandArgumentType
	Required        bool
	RequiredMessage string
	ValidValues     []interface{}
	ValidValuesFunc func() []interface{}
}

func (ca *CommandArgument) GetValidValues() []interface{} {
	if ca.ValidValues != nil {
		return ca.ValidValues
	}
	if ca.ValidValuesFunc != nil {
		return ca.ValidValuesFunc()
	}
	return nil
}

type Command struct {
	Name, Description string
	Aliases           []string
	Handler           CommandHandler
	Permission        *CommandPermission
	Arguments         []*CommandArgument
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

func (ctx *CommandContext) Embed(embed *discordgo.MessageEmbed) {
	ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Reference: ctx.Message.Reference(),
		Embed:     embed,
	})
}

func (ctx *CommandContext) SuccesEmbed(embed *discordgo.MessageEmbed) {
	embed.Color = 0x50fa7b
	ctx.Embed(embed)
}

func (ctx *CommandContext) ErrorEmbed(embed *discordgo.MessageEmbed) {
	embed.Color = 0xff5555
	ctx.Embed(embed)
}
