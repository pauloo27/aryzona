package fun

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/providers/livescore"
)

var UnFollowCommand = command.Command{
	Name:        "unfollow",
	Description: "Stop following a soccer match",
	Parameters: []*command.CommandParameter{
		{
			Name:        "game",
			Required:    false,
			Description: "team name or match id, if none is provided, all matches are unfollowed",
			Type:        parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		authorID := ctx.AuthorID

		if len(ctx.Args) == 0 {
			if len(userFollowedMatcheIDs[authorID]) == 0 {
				ctx.Error("You are not following any match")
				return
			}
			for _, matchID := range userFollowedMatcheIDs[authorID] {
				liveMatch, err := livescore.GetLiveMatch(matchID)
				if err != nil {
					ctx.Error(fmt.Sprintf("Something went wrong: %v", err))
					return
				}
				_ = liveMatch.RemoveListener(getListenerID(authorID, matchID))
				removeUserFollow(authorID, matchID)
			}
			ctx.Success("Unfollowed all matches")
		} else {
			teamNameOrID := ctx.Args[0].(string)
			match, err := getMatchByTeamNameOrID(teamNameOrID)
			if err != nil {
				ctx.Error(fmt.Sprintf("Something went wrong: %v", err))
				return
			}
			if match == nil {
				ctx.Error("Match not found")
				return
			}
			liveMatch, err := livescore.GetLiveMatch(match.ID)
			if err != nil {
				ctx.Error(fmt.Sprintf("Something went wrong: %v", err))
				return
			}
			err = liveMatch.RemoveListener(getListenerID(authorID, match.ID))
			if err == livescore.ErrListenerNotFound {
				ctx.Error("You are not following this match")
				return
			}
			removeUserFollow(authorID, match.ID)
			ctx.Success(fmt.Sprintf("Unfollowed match %s x %s", match.T1.Name, match.T2.Name))
		}
	},
}
