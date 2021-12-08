package utils

import (
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/providers/livescore"
	"github.com/Pauloo27/aryzona/utils"
)

var ScoreCommand = command.Command{
	Name: "score", Description: "Get a game score",
	Aliases: []string{"placar", "gols"},
	Arguments: []*command.CommandArgument{
		{
			Name:            "game",
			Required:        true,
			RequiredMessage: "Missing the team name or a match id",
			Description:     "team name or match id",
			Type:            command.ArgumentString,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		var match *livescore.MatchInfo
		teamNameOrID := ctx.Args[0].(string)
		if _, err := strconv.Atoi(teamNameOrID); err == nil {
			match, err = livescore.FetchMatchInfo(teamNameOrID)
			if err != nil {
				ctx.Error(err.Error())
				return
			}
		} else {
			match, err = livescore.FetchMatchInfoByTeamName(strings.Join(ctx.RawArgs, " "))
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

		desc := strings.Builder{}
		if len(match.Events) > 0 {
			for _, event := range match.Events {
				text := event.Text

				desc.WriteString(utils.Fmt(" -> %d' %s\n", event.Min, text))
			}
		}
		ctx.Embed(
			discord.NewEmbed().
				WithDescription(utils.Fmt("%s: %s, %s", match.CupName, match.StadiumName, match.StadiumCity)).
				WithColor(color).
				WithField("Time", match.Time).
				WithFieldInline(match.T1.Name, strconv.Itoa(match.T1.Score)).
				WithFieldInline(match.T2.Name, strconv.Itoa(match.T2.Score)).
				WithField("ID", match.ID).
				WithDescription(desc.String()),
		)
	},
}
