package bot

import (
	"fmt"
	"strings"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
	k "github.com/pauloo27/toolkit"
	"github.com/pauloo27/toolkit/slices"
)

var HelpCommand = command.Command{
	Name:    "help",
	Aliases: []string{"h"},
	Parameters: []*command.Parameter{
		{
			Name:     "command",
			Required: false, Type: parameters.ParameterString,
		},
		{
			Name:     "subcommand",
			Required: false, Type: parameters.ParameterString,
		},
	},
	Handler: func(ctx *command.Context) command.Result {
		if len(ctx.Args) == 0 {
			return listCommands(ctx)
		} else {
			return helpForCommand(ctx)
		}
	},
}

func listCommands(ctx *command.Context) command.Result {
	t := ctx.T.(*i18n.CommandHelp)

	embed := model.NewEmbed()
	sb := strings.Builder{}
	embed.WithTitle(t.Title.Str())
	lastCategory := ""
	for _, cmd := range command.GetCommandList() {
		cmdLang := i18n.MustGetCommandDefinition(ctx.Lang, cmd.Name)

		categoryName := cmd.GetCategory().Name
		localizedCategoryName := t.Categories[categoryName].Str()
		if localizedCategoryName == "" {
			logger.Warnf("Missing category name for %s", categoryName)
			localizedCategoryName = categoryName
		}

		if lastCategory != categoryName {
			sb.WriteString(fmt.Sprintf("\n**%s %s**:\n", cmd.GetCategory().Emoji, localizedCategoryName))
		}
		var permission string
		if cmd.Permission != nil {
			permission = t.RequiresPermission.Str(cmd.Permission.Name)
		}
		var aliases string
		if len(cmd.Aliases) > 0 {
			aliases = t.AKA.Str(strings.Join(cmd.Aliases, ", "))
		}
		sb.WriteString(fmt.Sprintf(
			" -> `%s%s` %s: **%s** %s\n",
			command.Prefix, cmdLang.Name, aliases, cmdLang.Description, permission,
		))
		lastCategory = cmd.GetCategory().Name
	}
	embed.WithFooter(
		t.ForMoreInfo.Str(
			command.Prefix, ctx.UsedName,
		),
	)

	embed.WithDescription(sb.String())

	return ctx.SuccessEmbed(embed)
}

func helpForCommand(ctx *command.Context) command.Result {
	t := ctx.T.(*i18n.CommandHelp)

	commandName := ctx.Args[0].(string)
	var subCommandName string
	if len(ctx.Args) > 1 {
		subCommandName = ctx.Args[1].(string)
	}

	cmd, found := command.GetCommandMap()[commandName]
	if !found {
		return ctx.Error(t.CommandNotFound.Str())
	}
	rootCmd := cmd

	if subCommandName != "" {
		subCommand := slices.Find(cmd.SubCommands, func(subCmd *command.Command) bool {
			return subCmd.Name == subCommandName
		})
		if subCommand == nil {
			return ctx.Error(t.SubCommandNotFound.Str())
		}
		cmd = *subCommand
	}

	fullCommandName := rootCmd.Name
	if rootCmd != cmd {
		fullCommandName = fmt.Sprintf("%s %s", rootCmd.Name, cmd.Name)
	}

	cmdLang := i18n.MustGetCommandDefinition(ctx.Lang, cmd.Name)

	categoryName := cmd.GetCategory().Name
	localizedCategoryName := t.Categories[categoryName].Str()

	embed := model.NewEmbed().
		WithTitle(fullCommandName).
		WithField(t.Category.Str(), fmt.Sprintf("%s %s", rootCmd.GetCategory().Emoji, localizedCategoryName)).
		WithDescription(cmdLang.Description.Str())

	if cmd.Aliases != nil {
		embed.WithField(t.Aliases.Str(), strings.Join(cmd.Aliases, ", "))
	}

	if cmd.Permission != nil {
		embed.WithField(t.Permission.Str(), cmd.Permission.Name)
	}

	if cmd.SubCommands != nil {
		embed.WithField(
			t.SubCommands.Str(),
			strings.Join(
				slices.Map(cmd.SubCommands, func(cmd *command.Command) string {
					return cmd.Name
				}),
				", ",
			),
		)
	}

	if cmd.Validations != nil {
		embed.WithField(
			t.Validations.Str(),
			strings.Join(
				slices.Map(cmd.Validations, func(validation *command.Validation) string {
					v, err := t.RawMap.Get("common", "validations", validation.Name, "description")
					if err != nil {
						logger.Errorf("Missing validation description for %s", validation.Name)
						return validation.Name
					}
					return v.(string)
				}),
				", ",
			),
		)
	}

	if cmd.Parameters != nil {
		i := 0
		embed.WithField(
			t.Parameters.Str(),
			strings.Join(
				slices.Map(cmd.Parameters, func(param *command.Parameter) string {
					paramLang := cmdLang.Parameters[i]
					i++
					return fmt.Sprintf("%s: %s (%s)",
						paramLang.Name, paramLang.Description.Str(),
						k.Is(param.Required, t.Required.Str(), t.NotRequired.Str()))
				}),
				"\n",
			),
		)
	}

	return ctx.SuccessEmbed(
		embed,
	)
}
