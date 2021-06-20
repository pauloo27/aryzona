package utils

import (
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/providers/livescore"
	"github.com/Pauloo27/aryzona/utils"
)

var ScoreCommand = command.Command{
	Name: "score", Description: "Get a game score",
	Aliases: []string{"placar", "gols"},
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 0 {
			ctx.Error("Missing game id or team name")
			return
		}

		var match *livescore.MatchInfo
		if _, err := strconv.Atoi(ctx.Args[0]); err == nil {
			match, err = livescore.FetchMatchInfo(ctx.Args[0])
			if err != nil {
				ctx.Error(err.Error())
				return
			}
		} else {
			match, err = livescore.FetchMatchInfoByTeamName(strings.Join(ctx.Args, " "))
			if err != nil {
				ctx.Error(err.Error())
				return
			}
			if match == nil {
				ctx.Error("Match not found")
				return
			}
		}

		var color int
		if match.T1.Score == match.T2.Score {
			color = 0xC0FFEE
		} else if match.T1.Score > match.T2.Score {
			color = match.T1.ColorAsInt()
		} else {
			color = match.T2.ColorAsInt()
		}
		ctx.Embed(
			utils.NewEmbedBuilder().
				Description(utils.Fmt("%s: %s, %s", match.CupName, match.StadiumName, match.StadiumCity)).
				Color(color).
				Field("Time", match.Time).
				FieldInline(match.T1.Name, strconv.Itoa(match.T1.Score)).
				FieldInline(match.T2.Name, strconv.Itoa(match.T2.Score)).
				Build(),
		)
	},
}
