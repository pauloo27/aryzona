package dcgo

import (
	"fmt"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

var discordTypeMap = map[*command.CommandArgumentType]discordgo.ApplicationCommandOptionType{
	command.ArgumentString: discordgo.ApplicationCommandOptionString,
	command.ArgumentText:   discordgo.ApplicationCommandOptionString,
	command.ArgumentInt:    discordgo.ApplicationCommandOptionInteger,
	command.ArgumentBool:   discordgo.ApplicationCommandOptionBoolean,
}

func registerCommands(bot DcgoBot) error {
	session := bot.d.s

	mustGetChoisesFor := func(arg *command.CommandArgument) (options []*discordgo.ApplicationCommandOptionChoice) {
		for _, value := range arg.GetValidValues() {
			options = append(options, &discordgo.ApplicationCommandOptionChoice{
				Name:  fmt.Sprintf("%v", value),
				Value: value,
			})
		}
		return
	}

	mustGetTypeFor := func(arg *command.CommandArgument) discordgo.ApplicationCommandOptionType {
		t, found := discordTypeMap[arg.Type]
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

		for _, arg := range cmd.Arguments {
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
			args = append(args, fmt.Sprintf("%v", option.Value))
		}

		var authorID string

		if i.Member == nil {
			authorID = i.User.ID
		} else {
			authorID = i.Member.User.ID
		}

		event := command.Event{
			AuthorID: authorID,
			GuildID:  i.GuildID,
			Reply: func(message string) error {
				return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: message,
					},
				})
			},
			ReplyEmbed: func(embed *discord.Embed) error {
				return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{},
					},
				})
			},
		}
		command.HandleCommand(commandName, args, &event, bot)
	})

	return nil
}
