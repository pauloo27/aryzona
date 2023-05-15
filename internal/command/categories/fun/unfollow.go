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
	Parameters: []*command.CommandParameter{
		{
			Name:     "game",
			Required: false,
			Type:     parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandUnFollow)

		authorID := ctx.AuthorID

		if len(ctx.Args) == 0 {
			if len(followedMatchIDs[authorID]) == 0 {
				ctx.Error(t.NotFollowingAny.Str())
				return
			}
			for _, matchID := range followedMatchIDs[authorID] {
				liveMatch, err := livescore.GetLiveMatch(matchID)
				if err != nil {
					ctx.Error(t.SomethingWentWrong.Str())
					return
				}
				_ = liveMatch.RemoveListener(getListenerID(authorID, matchID))
				removeUserFollow(authorID, matchID)
			}
			ctx.Success(t.UnFollowedAll.Str())
		} else {
			teamNameOrID := ctx.Args[0].(string)
			match, err := getMatchByTeamNameOrID(teamNameOrID)
			if err != nil {
				ctx.Error(t.SomethingWentWrong.Str())
				logger.Error(err)
				return
			}
			if match == nil {
				ctx.Error(t.MatchNotFound.Str())
				return
			}
			liveMatch, err := livescore.GetLiveMatch(match.ID)
			if err != nil {
				ctx.Error(t.SomethingWentWrong.Str())
				logger.Error(err)
				return
			}
			err = liveMatch.RemoveListener(getListenerID(authorID, match.ID))
			if err == livescore.ErrListenerNotFound {
				ctx.Error(t.NotFollowingMatch.Str())
				return
			}
			removeUserFollow(authorID, match.ID)
			ctx.Success(t.UnfollowedMatch.Str(match.T1.Name, match.T2.Name))
		}
	},
}
