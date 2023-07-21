package fun

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/livescore"
	"github.com/pauloo27/logger"
)

var UnFollowCommand = command.Command{
	Name: "unfollow",
	Parameters: []*command.Parameter{
		{
			Name:     "game",
			Required: false,
			Type:     parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandUnFollow)

		authorID := ctx.AuthorID

		if len(ctx.Args) == 0 {
			if len(followedMatchIDs[authorID]) == 0 {
				return ctx.Error(t.NotFollowingAny.Str())
			}
			for _, matchID := range followedMatchIDs[authorID] {
				liveMatch, err := livescore.GetLiveMatch(matchID)
				if err != nil {
					return ctx.Error(t.SomethingWentWrong.Str())
				}
				_ = liveMatch.RemoveListener(getListenerID(authorID, matchID))
				removeUserFollow(authorID, matchID)
			}
			return ctx.Success(t.UnFollowedAll.Str())
		}

		teamNameOrID := ctx.Args[0].(string)
		match, err := getMatchByTeamNameOrID(teamNameOrID)
		if err != nil {
			logger.Error(err)
			return ctx.Error(t.SomethingWentWrong.Str())
		}
		if match == nil {
			return ctx.Error(t.MatchNotFound.Str())
		}
		liveMatch, err := livescore.GetLiveMatch(match.ID)
		if err != nil {
			logger.Error(err)
			return ctx.Error(t.SomethingWentWrong.Str())
		}
		err = liveMatch.RemoveListener(getListenerID(authorID, match.ID))
		if err == livescore.ErrListenerNotFound {
			return ctx.Error(t.NotFollowingMatch.Str())
		}
		removeUserFollow(authorID, match.ID)
		return ctx.Success(t.UnfollowedMatch.Str(match.T1.Name, match.T2.Name))
	},
}
