package fun

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/aryzona/internal/providers/livescore"
	"github.com/Pauloo27/logger"
)

const (
	maxFollowPerUser = 3
)

var (
	userFollowedMatcheIDs = make(map[string][]string)
)

var FollowCommand = command.Command{
	Name:    "follow",
	Aliases: []string{"live"},
	Parameters: []*command.CommandParameter{
		{
			Name:            "game",
			Required:        true,
			RequiredMessage: "Missing the team name or a match id",
			Type:            parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandFollow)

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
		if errors.Is(err, livescore.ErrMatchHasFinished) {
			ctx.Error(t.MatchFinished.Str())
			return
		}

		if len(userFollowedMatcheIDs[ctx.AuthorID]) >= maxFollowPerUser {
			ctx.Error(t.FollowLimitReached.Str(maxFollowPerUser, command.Prefix))
			return
		}

		for _, followedMatchID := range userFollowedMatcheIDs[ctx.AuthorID] {
			if followedMatchID == match.ID {
				ctx.Error(t.AlreadyFollowing.Str())
				return
			}
		}

		matchInfoI18n := &MatchInfoI18n{
			Match:       t.Match.Str(),
			Time:        t.Time.Str(),
			TimePenalty: t.TimePenalty.Str(),
		}

		addUserFollow(ctx.AuthorID, match.ID)

		embed := buildMatchEmbed(match, matchInfoI18n)
		ctx.Embed(embed)

		listenerID := getListenerID(ctx.AuthorID, match.ID)

		liveMatch.AddListener(listenerID, func(match *livescore.LiveMatch, err error) {
			if err != nil {
				ctx.Error(err.Error())
				removeUserFollow(ctx.AuthorID, match.MatchID)
				return
			}
			embed := buildMatchEmbed(match.CurrentData, matchInfoI18n)
			err = ctx.EditEmbed(embed)
			if err != nil {
				_ = liveMatch.RemoveListener(listenerID)
				removeUserFollow(ctx.AuthorID, match.MatchID)
			}
		})
	},
}

func addUserFollow(userID string, matchID string) {
	userFollowedMatcheIDs[userID] = append(userFollowedMatcheIDs[userID], matchID)
}

func removeUserFollow(userID string, matchID string) {
	for i, id := range userFollowedMatcheIDs[userID] {
		if id == matchID {
			userFollowedMatcheIDs[userID] = append(userFollowedMatcheIDs[userID][:i], userFollowedMatcheIDs[userID][i+1:]...)
		}
	}
}

func getListenerID(authorID string, matchID string) string {
	return fmt.Sprintf("%s-%s", authorID, matchID)
}

func getMatchByTeamNameOrID(teamNameOrID string) (*livescore.MatchInfo, error) {
	if _, err := strconv.Atoi(teamNameOrID); err == nil {
		return livescore.FetchMatchInfo(teamNameOrID)
	}
	return livescore.FetchMatchInfoByTeamName(teamNameOrID)
}
