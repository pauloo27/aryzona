package fun

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/providers/livescore"
	"github.com/Pauloo27/aryzona/internal/utils"
)

const (
	maxFollowPerUser = 3
)

var (
	// TODO: make sure something will remove stuff from memory
	userFollows = make(map[string]int)
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

		if match.Time == "FT" {
			ctx.Error("Match already finished")
			return
		}

		if userFollows[ctx.AuthorID] >= maxFollowPerUser {
			ctx.Error(
				fmt.Sprintf(
					"You cannot follow more than 3 matches, you can `%sunfollow` some match to follow another one",
					command.Prefix,
				),
			)
			return
		}

		userFollows[ctx.AuthorID]++

		embed := BuildMatchEmbed(match)
		err := ctx.ReplyEmbed(embed)
		if err != nil {
			// TODO: cancel fetcher? (only if the only one using it)
			ctx.Error("Something went wrong")
			return
		}

		utils.Go(func() {
			for {
				time.Sleep(60 * time.Second)
				match, err := livescore.FetchMatchInfo(match.ID)
				if err != nil {
					ctx.Error(err.Error())
					return
				}
				if match.Time == "FT" {
					userFollows[ctx.AuthorID]--
					if userFollows[ctx.AuthorID] <= 0 {
						delete(userFollows, ctx.AuthorID)
					}
					return
				}
				embed := BuildMatchEmbed(match)
				err = ctx.EditEmbed(embed)
				if err != nil {
					// TODO: cancel fetcher? (only if the only one using it)
					ctx.Error("Something went wrong")
					return
				}
			}
		})
	},
}
