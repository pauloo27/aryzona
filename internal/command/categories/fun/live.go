package fun

import (
	"errors"
	"strconv"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/providers/livescore"
)

const (
	maxFollowPerUser = 3
)

var LiveCommand = command.Command{
	Name:        "live",
	Aliases:     []string{"follow"},
	Description: "Get live updates from a soccer match",
	Parameters: []*command.CommandParameter{
		{
			Name:            "game",
			Required:        true,
			RequiredMessage: "Missing the team name or a match id",
			Description:     "team name or match id",
			Type:            parameters.ParameterText,
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
			match, err = livescore.FetchMatchInfoByTeamName(teamNameOrID)
			if err != nil {
				ctx.Error(err.Error())
				return
			}
			if match == nil {
				ctx.Error("Match not found")
				return
			}
		}

		liveMatch, err := livescore.GetLiveMatch(match.ID)
		if errors.Is(err, livescore.ErrMatchHasFinished) {
			ctx.Error("Match has finished")
			return
		}

		embed := BuildMatchEmbed(match)
		ctx.Embed(embed)

		liveMatch.AddListener(func(match *livescore.LiveMatch, err error) {
			if err != nil {
				ctx.Error(err.Error())
				return
			}
			embed := BuildMatchEmbed(match.CurrentData)
			// TODO: remove handler
			_ = ctx.EditEmbed(embed)
		})
	},
}
