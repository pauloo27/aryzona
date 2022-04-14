package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/providers/livescore"
)

var ScoreCommand = command.Command{
	Name: "score", Description: "Show live matches score",
	Aliases: []string{"placar", "gols", "live", "scores"},
	Parameters: []*command.CommandParameter{
		{
			Name:            "game",
			Required:        false,
			RequiredMessage: "Missing the team name or a match id",
			Description:     "team name or match id",
			Type:            parameters.ParameterString,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 1 {
			ShowMatchInfo(ctx)
			return
		}
		ListLiveMatches(ctx)
	},
}

func ListLiveMatches(ctx *command.CommandContext) {
	matches, err := livescore.ListLives()
	if err != nil {
		ctx.Error(err.Error())
		return
	}
	if len(matches) == 0 {
		ctx.Error("I didn't find any live match right now...")
		return
	}
	desc := strings.Builder{}
	for _, match := range matches {
		desc.WriteString(fmt.Sprintf("**%s** ||(%d) x (%d)|| **%s**: _%s_\n",
			match.T1.Name, match.T1.Score,
			match.T2.Score, match.T2.Name,
			match.ID,
		))
	}
	ctx.SuccessEmbed(
		discord.NewEmbed().
			WithTitle("⚽ Live matches:").
			WithFooter(
				fmt.Sprintf("Use `%s%s <match id or team name>` to see details",
					command.Prefix, ctx.UsedName,
				)).
			WithDescription(desc.String()),
	)
}

func ShowMatchInfo(ctx *command.CommandContext) {
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

			desc.WriteString(fmt.Sprintf(" -> %d' %s\n", event.Min, text))
		}
	}

	descStr := desc.String()
	if len(descStr) > 4096 {
		descStr = descStr[:4093] + "..."
	}

	ctx.Embed(
		discord.NewEmbed().
			WithDescription(fmt.Sprintf("%s: %s, %s", match.CupName, match.StadiumName, match.StadiumCity)).
			WithColor(color).
			WithField("Time", match.Time).
			WithFieldInline(match.T1.Name, strconv.Itoa(match.T1.Score)).
			WithFieldInline(match.T2.Name, strconv.Itoa(match.T2.Score)).
			WithField("ID", match.ID).
			WithDescription(descStr),
	)
}
