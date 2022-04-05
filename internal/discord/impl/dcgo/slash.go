package dcgo

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

var discordTypeMap = map[*command.BaseType]discordgo.ApplicationCommandOptionType{
	parameters.TypeString: discordgo.ApplicationCommandOptionString,
	parameters.TypeInt:    discordgo.ApplicationCommandOptionInteger,
	parameters.TypeBool:   discordgo.ApplicationCommandOptionBoolean,
}

func registerCommands(bot DcgoBot) error {
	session := bot.d.s

	mustGetChoisesFor := func(arg *command.CommandParameter) (options []*discordgo.ApplicationCommandOptionChoice) {
		for _, value := range arg.GetValidValues() {
			options = append(options, &discordgo.ApplicationCommandOptionChoice{
				Name:  fmt.Sprintf("%v", value),
				Value: value,
			})
		}
		return
	}

	mustGetTypeFor := func(arg *command.CommandParameter) discordgo.ApplicationCommandOptionType {
		t, found := discordTypeMap[arg.Type.BaseType]
		if !found {
			logger.Fatalf("cannot find discord type for %s", arg.Type.Name)
		}
		return t
	}

	var slashCommands []*discordgo.ApplicationCommand

	for key, cmd := range command.GetCommandMap() {

		// skip aliases
		if key != cmd.Name {
			continue
		}

		slashCommand := discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		}

		if cmd.SubCommands != nil {
			for _, subCmd := range cmd.SubCommands {
				slashCommandOptions := []*discordgo.ApplicationCommandOption{}
				for _, arg := range subCmd.Parameters {
					slashCommandOptions = append(slashCommandOptions, &discordgo.ApplicationCommandOption{
						Name:        arg.Name,
						Description: arg.Description,
						Type:        mustGetTypeFor(arg),
						Choices:     mustGetChoisesFor(arg),
						Required:    arg.Required,
					})
				}
				slashCommand.Options = append(slashCommand.Options, &discordgo.ApplicationCommandOption{
					Name:        subCmd.Name,
					Description: subCmd.Description,
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options:     slashCommandOptions,
				})
			}
		}

		for _, arg := range cmd.Parameters {
			slashCommand.Options = append(slashCommand.Options, &discordgo.ApplicationCommandOption{
				Name:        arg.Name,
				Description: arg.Description,
				Required:    arg.Required,
				Type:        mustGetTypeFor(arg),
				Choices:     mustGetChoisesFor(arg),
			})
		}

		slashCommands = append(slashCommands, &slashCommand)
	}
	_, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", slashCommands)
	if err != nil {
		return err
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandName := i.ApplicationCommandData().Name
		_, ok := command.GetCommandMap()[commandName]
		if !ok {
			logger.Error("Invalid slash command interaction received:", i.ApplicationCommandData().Name)
			return
		}

		var args []string
		for _, option := range i.ApplicationCommandData().Options {
			if option.Type == discordgo.ApplicationCommandOptionSubCommand {
				args = append(args, option.Name)
				for _, subCommandOption := range option.Options {
					args = append(args, fmt.Sprintf("%v", subCommandOption.Value))
				}
				break
			}
			args = append(args, fmt.Sprintf("%v", option.Value))
		}

		edit := func(msg string, embed *discord.Embed) error {
			var embeds []*discordgo.MessageEmbed
			if embed != nil {
				embeds = append(embeds, buildEmbed(embed))
			}
			_, err := s.InteractionResponseEdit(s.State.User.ID,
				i.Interaction,
				&discordgo.WebhookEdit{
					Embeds:  embeds,
					Content: msg,
				},
			)
			return err
		}

		respond := func(msg string, embed *discord.Embed) error {
			var embeds []*discordgo.MessageEmbed
			if embed != nil {
				embeds = append(embeds, buildEmbed(embed))
			}
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds:  embeds,
					Content: msg,
				},
			})
		}

		var authorID string

		if i.Member == nil {
			authorID = i.User.ID
		} else {
			authorID = i.Member.User.ID
		}

		event := command.Adapter{
			AuthorID: authorID,
			GuildID:  i.GuildID,
			DeferResponse: func() error {
				return s.InteractionRespond(
					i.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					},
				)
			},
			Reply: func(ctx *command.CommandContext, msg string) error {
				if ctx.Command.Deferred {
					return edit(msg, nil)
				}
				return respond(msg, nil)
			},
			ReplyEmbed: func(ctx *command.CommandContext, embed *discord.Embed) error {
				if ctx.Command.Deferred {
					return edit("", embed)
				}
				return respond("", embed)
			},
		}
		command.HandleCommand(commandName, args, &event, bot)
	})

	return nil
}
