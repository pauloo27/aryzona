package utils

import (
	"strings"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/providers/livescore"
	"github.com/Pauloo27/aryzona/utils"
)

var LiveCommand = command.Command{
	Name: "live", Description: "List live games",
	Aliases: []string{"matches", "jogos"},
	Handler: func(ctx *command.CommandContext) {
		matches, err := livescore.ListLives()
		if err != nil {
			ctx.Error(err.Error())
			return
		}
		if len(matches) == 0 {
			ctx.Error("I didn't find any live match right now...")
		}
		// more than 25? cant use fields =(
		desc := strings.Builder{}
		for _, match := range matches {
			desc.WriteString(utils.Fmt("**%s** (%d) x (%d) **%s**: _%s_\n",
				match.T1.Name, match.T1.Score,
				match.T2.Score, match.T2.Name,
				match.Id,
			))
		}
		ctx.SuccesEmbed(
			utils.NewEmbedBuilder().
				Title("Live matches:").
				Description(desc.String()).
				Build(),
		)
	},
}
