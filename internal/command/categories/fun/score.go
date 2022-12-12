package fun

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
	Name: "score", Description: "Show matches score",
	Aliases: []string{"placar", "gols", "scores"},
	Parameters: []*command.CommandParameter{
		{
			Name:            "game",
			Required:        false,
			RequiredMessage: "Missing the team name or a match id",
			Description:     "team name or match id",
			Type:            parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 1 {
			showMatchInfo(ctx)
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
		desc.WriteString(fmt.Sprintf("%s **%s** ||(%d) x (%d)|| **%s**\n",
			match.Time,
			match.T1.Name, match.T1Score,
			match.T2Score, match.T2.Name,
		))
	}
	ctx.SuccessEmbed(
		discord.NewEmbed().
			WithTitle("⚽ Live matches:").
			WithFooter(
				fmt.Sprintf("Use `%s%s <team name>` to see details",
					command.Prefix, ctx.UsedName,
				)).
			WithDescription(desc.String()),
	)
}

func showMatchInfo(ctx *command.CommandContext) {
	teamNameOrID := ctx.Args[0].(string)
	match, err := getMatchByTeamNameOrID(teamNameOrID)

	if err != nil {
		ctx.Error(err.Error())
		return
	}
	if match == nil {
		ctx.Error("Match not found")
		return
	}

	embed := buildMatchEmbed(match).
		WithFooter(fmt.Sprintf("Use `%slive %s` to get live updates", command.Prefix, teamNameOrID))

	ctx.Embed(embed)
}

func buildMatchEmbed(match *livescore.MatchInfo) *discord.Embed {
	desc := strings.Builder{}

	if len(match.Events) > 0 {
		for _, event := range match.Events {
			prefix, found := eventTypePrefixes[event.Type]
			if !found {
				continue
			}

			var eventTime string
			if event.Half == 4 {
				eventTime = "Pen"
			} else if event.ExtraMinute != 0 {
				eventTime += fmt.Sprintf("%d+%d'", event.Minute, event.ExtraMinute)
			} else {
				eventTime += fmt.Sprintf("%d'", event.Minute)
			}

			desc.WriteString(fmt.Sprintf(" -> %s %s [%s] %s\n", eventTime, prefix, event.Team.Name, event.PlayerName))
		}
	}

	descStr := desc.String()
	if len(descStr) > 4096 {
		descStr = descStr[:4093] + "..."
	}

	var t1Score, t2Score string

	if match.T1Score != -1 {
		t1Score = strconv.Itoa(match.T1Score)
	} else {
		t1Score = "_"
	}

	if match.T2Score != -1 {
		t2Score = strconv.Itoa(match.T2Score)
	} else {
		t2Score = "_"
	}

	return discord.NewEmbed().
		WithColor(0xC0FFEE).
		WithField("Match", fmt.Sprintf("%s: %s, %s", match.CupName, match.StadiumName, match.StadiumCity)).
		WithField("Time", match.Time).
		WithImage(match.GetBannerURL()).
		WithFieldInline(match.T1.Name, t1Score).
		WithFieldInline(match.T2.Name, t2Score).
		WithDescription(descStr)
}

var eventTypePrefixes = map[livescore.EventType]string{
	livescore.EventTypeYellowCard:       "🟡",
	livescore.EventTypeDoubleYellowCard: "🟡+🟡=🔴",
	livescore.EventTypeRedCard:          "🔴",
	livescore.EventTypeGoal:             "⚽",
	livescore.EventTypeOvertimeGoal:     "⚽",
	livescore.EventTypeFoulPenaltyGoal:  "⚽",
	livescore.EventTypePenaltyGoal:      "✅",
	livescore.EventTypePenaltyMissed:    "❌",
}
