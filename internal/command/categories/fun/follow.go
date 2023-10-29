package fun

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/livescore"
	"github.com/pauloo27/logger"
)

const (
	maxFollowPerUser = 3
)

var (
	followedMatchIDs = make(map[string][]string)
)

var FollowCommand = command.Command{
	Name: "follow",
	Parameters: []*command.Parameter{
		{
			Name:     "game",
			Required: true,
			Type:     parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandFollow)

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
		if errors.Is(err, livescore.ErrMatchHasFinished) {
			return ctx.Error(t.MatchFinished.Str())
		}

		if len(followedMatchIDs[ctx.AuthorID]) >= maxFollowPerUser {
			return ctx.Error(t.FollowLimitReached.Str(maxFollowPerUser, command.Prefix))
		}

		for _, followedMatchID := range followedMatchIDs[ctx.AuthorID] {
			if followedMatchID == match.ID {
				return ctx.Error(t.AlreadyFollowing.Str())
			}
		}

		addUserFollow(ctx.AuthorID, match.ID)

		embed := buildMatchEmbed(match, t.MatchInfo)

		listenerID := getListenerID(ctx.AuthorID, match.ID)

		liveMatch.AddListener(listenerID, func(match *livescore.LiveMatch, err error) {
			if err != nil {
				ctx.Error(err.Error())
				removeUserFollow(ctx.AuthorID, match.MatchID)
				return
			}
			embed := buildMatchEmbed(match.CurrentData, t.MatchInfo)
			err = ctx.EditComplexMessage(&model.ComplexMessage{
				Embeds: []*model.Embed{embed},
			})
			if err != nil {
				_ = liveMatch.RemoveListener(listenerID)
				removeUserFollow(ctx.AuthorID, match.MatchID)
			}
		})

		return ctx.Embed(embed)
	},
}

func addUserFollow(userID string, matchID string) {
	followedMatchIDs[userID] = append(followedMatchIDs[userID], matchID)
}

func removeUserFollow(userID string, matchID string) {
	for i, id := range followedMatchIDs[userID] {
		if id == matchID {
			followedMatchIDs[userID] = append(
				followedMatchIDs[userID][:i],
				followedMatchIDs[userID][i+1:]...,
			)
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
