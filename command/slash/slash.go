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
		event := command.Event{
			Message: i.Message,
		}
		command.HandleCommand(commandName, args, s, &event)
	})

	return nil
}
