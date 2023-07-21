package fun

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/livescore"
)

var ScoreCommand = command.Command{
	Name: "score",
	Parameters: []*command.Parameter{
		{
			Name:     "game",
			Required: false,
			Type:     parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandScore)

		if len(ctx.Args) == 1 {
			return showMatchInfo(ctx, t)
		}

		return ListLiveMatches(ctx, t)
	},
}

func ListLiveMatches(ctx *command.Context, t *i18n.CommandScore) command.Result {
	matches, err := livescore.ListLives()
	if err != nil {
		return ctx.Error(err.Error())
	}
	if len(matches) == 0 {
		return ctx.Error(t.NoMatchesLive.Str())
	}
	desc := strings.Builder{}
	for _, match := range matches {
		desc.WriteString(fmt.Sprintf("%s **%s** ||(%d) x (%d)|| **%s**\n",
			match.Time,
			match.T1.Name, match.T1Score,
			match.T2Score, match.T2.Name,
		))
	}

	return ctx.SuccessEmbed(
		model.NewEmbed().
			WithTitle(t.Title.Str(":soccer:")).
			WithFooter(t.Footer.Str(command.Prefix, ctx.UsedName)).
			WithDescription(desc.String()),
	)
}

func showMatchInfo(ctx *command.Context, t *i18n.CommandScore) command.Result {
	teamNameOrID := ctx.Args[0].(string)
	match, err := getMatchByTeamNameOrID(teamNameOrID)

	if err != nil {
		return ctx.Error(err.Error())
	}
	if match == nil {
		return ctx.Error(t.MatchNotFound.Str())
	}

	embed := buildMatchEmbed(match, t.MatchInfo).
		WithFooter(
			t.LiveUpdates.Str(command.Prefix, teamNameOrID),
		)

	return ctx.Embed(embed)
}

func buildMatchEmbed(match *livescore.MatchInfo, t *i18n.MatchInfo) *model.Embed {
	desc := strings.Builder{}

	if len(match.Events) > 0 {
		for _, event := range match.Events {
			prefix, found := eventTypePrefixes[event.Type]
			if !found {
				continue
			}

			var eventTime string
			if event.Half == 4 {
				eventTime = t.TimePenalty.Str()
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

	return model.NewEmbed().
		WithColor(0xC0FFEE).
		WithField(t.Match.Str(), fmt.Sprintf("%s: %s, %s", match.CupName, match.StadiumName, match.StadiumCity)).
		WithField(t.Time.Str(), match.Time).
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
