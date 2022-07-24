package bot

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/aryzona/internal/utils/slice"
)

var HelpCommand = command.Command{
	Name: "help", Description: "List all commands",
	Aliases: []string{"h"},
	Parameters: []*command.CommandParameter{
		{
			Name: "command", Description: "Command to get help for",
			Required: false, Type: parameters.ParameterString,
		},
		{
			Name: "subcommand", Description: "Sub command to get help for",
			Required: false, Type: parameters.ParameterString,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 0 {
			listCommands(ctx)
		} else {
			helpForCommand(ctx)
		}
	},
}

func listCommands(ctx *command.CommandContext) {
	embed := discord.NewEmbed()
	sb := strings.Builder{}
	embed.WithTitle("List of commands")
	lastCategory := ""
	for _, cmd := range command.GetCommandList() {
		if lastCategory != cmd.GetCategory().Name {
			sb.WriteString(fmt.Sprintf("\n**%s %s**:\n", cmd.GetCategory().Emoji, cmd.GetCategory().Name))
		}
		var permission string
		if cmd.Permission != nil {
			permission = fmt.Sprintf("(_requires you to... %s_)", cmd.Permission.Name)
		}
		var aliases string
		if len(cmd.Aliases) > 0 {
			aliases = fmt.Sprintf("(aka %s)", strings.Join(cmd.Aliases, ", "))
		}
		sb.WriteString(fmt.Sprintf(
			" - `%s%s` %s: **%s** %s\n",
			command.Prefix, cmd.Name, aliases, cmd.Description, permission,
		))
		lastCategory = cmd.GetCategory().Name
	}
	embed.WithFooter(
		fmt.Sprintf("For more information on a command, use `%s%s <command>`",
			command.Prefix, ctx.UsedName,
		),
	)
	ctx.SuccessEmbed(embed.WithDescription(sb.String()))
}

func helpForCommand(ctx *command.CommandContext) {
	commandName := ctx.Args[0].(string)
	var subCommandName string
	if len(ctx.Args) > 1 {
		subCommandName = ctx.Args[1].(string)
	}

	cmd, found := command.GetCommandMap()[commandName]
	if !found {
		ctx.Error("Command not found")
		return
	}
	rootCmd := cmd

	if subCommandName != "" {
		subCommand := slice.Find(cmd.SubCommands, func(subCmd *command.Command) bool {
			return subCmd.Name == subCommandName
		})
		if subCommand == nil {
			ctx.Error("Sub command not found")
			return
		}
		cmd = *subCommand
	}

	fullCommandName := rootCmd.Name
	if rootCmd != cmd {
		fullCommandName = fmt.Sprintf("%s %s", rootCmd.Name, cmd.Name)
	}

	embed := discord.NewEmbed().
		WithTitle(fullCommandName).
		WithField("Category", fmt.Sprintf("%s %s", rootCmd.GetCategory().Emoji, rootCmd.GetCategory().Name)).
		WithDescription(cmd.Description)

	if cmd.Aliases != nil {
		embed.WithField("Aliases", strings.Join(cmd.Aliases, ", "))
	}

	if cmd.Permission != nil {
		embed.WithField("Required Permission", cmd.Permission.Name)
	}

	if cmd.SubCommands != nil {
		embed.WithField(
			"Sub Commands",
			strings.Join(
				slice.Map(cmd.SubCommands, func(cmd *command.Command) string {
					return cmd.Name
				}),
				", ",
			),
		)
	}

	if cmd.Validations != nil {
		embed.WithField(
			"Validations",
			strings.Join(
				slice.Map(cmd.Validations, func(validation *command.CommandValidation) string {
					return validation.Description
				}),
				", ",
			),
		)
	}

	if cmd.Parameters != nil {
		embed.WithField(
			"Parameters",
			strings.Join(
				slice.Map(cmd.Parameters, func(param *command.CommandParameter) string {
					return fmt.Sprintf("%s: %s (%s)",
						param.Name, param.Description,
						utils.ConditionalString(param.Required, "required", "not required"))
				}),
				", ",
			),
		)
	}

	ctx.SuccessEmbed(
		embed,
	)
}
