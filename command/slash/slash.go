package slash

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

func RegisterCommands() error {
	for key, cmd := range command.GetCommandMap() {
		// skip aliases
		if key != cmd.Name {
			continue
		}

		slashCommand := discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		}

		_, err := discord.Session.ApplicationCommandCreate(discord.Session.State.User.ID, "", &slashCommand)
		if err != nil {
			return err
		}
	}

	discord.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandName := i.ApplicationCommandData().Name
		_, ok := command.GetCommandMap()[commandName]
		if !ok {
			logger.Error("Invalid slash command interaction received:", i.ApplicationCommandData().Name)
			return
		}

		var args []string
		for _, option := range i.ApplicationCommandData().Options {
			args = append(args, option.StringValue())
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
			ReplyEmbed: func(embed *discordgo.MessageEmbed) error {
				return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
			},
		}
		command.HandleCommand(commandName, args, s, &event)
	})

	return nil
}
