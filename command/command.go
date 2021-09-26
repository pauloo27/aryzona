package command

import (
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(*CommandContext)
type CommandPermissionChecker func(*CommandContext) bool

type CommandContext struct {
	Session           *discordgo.Session
	RawArgs           []string
	Args              []interface{}
	Reply             func(string) error
	ReplyEmbed        func(*discordgo.MessageEmbed) error
	AuthorID, GuildID string
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
	category          *CommandCategory
}

func (c *Command) GetCategory() *CommandCategory {
	return c.category
}

func (ctx *CommandContext) handleCannotSendMessage(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func (ctx *CommandContext) Success(message string) {
	ctx.handleCannotSendMessage(ctx.SuccessReturning(message))
}

func (ctx *CommandContext) SuccessReturning(message string) error {
	return ctx.Reply(utils.Fmt(":green_square: %s", message))
}

func (ctx *CommandContext) Error(message string) {
	ctx.handleCannotSendMessage(ctx.ErrorReturning(message))
}

func (ctx *CommandContext) ErrorReturning(message string) error {
	return ctx.Reply(utils.Fmt(":red_square: %s", message))
}

func (ctx *CommandContext) Embed(embed *discordgo.MessageEmbed) {
	ctx.handleCannotSendMessage(ctx.EmbedReturning(embed))
}

func (ctx *CommandContext) EmbedReturning(embed *discordgo.MessageEmbed) error {
	return ctx.ReplyEmbed(embed)
}

func (ctx *CommandContext) SuccessEmbed(embed *discordgo.MessageEmbed) {
	ctx.handleCannotSendMessage(ctx.SuccessEmbedReturning(embed))
}

func (ctx *CommandContext) SuccessEmbedReturning(embed *discordgo.MessageEmbed) error {
	embed.Color = 0x50fa7b
	return ctx.EmbedReturning(embed)
}

func (ctx *CommandContext) ErrorEmbed(embed *discordgo.MessageEmbed) {
	ctx.handleCannotSendMessage(ctx.ErrorEmbedReturning(embed))
}

func (ctx *CommandContext) ErrorEmbedReturning(embed *discordgo.MessageEmbed) error {
	embed.Color = 0xff5555
	return ctx.EmbedReturning(embed)
}
