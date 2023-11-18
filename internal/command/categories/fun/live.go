package fun

import (
	"fmt"
	"math"
	"strings"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/livescore"
)

const (
	PageSize = 20
)

var LiveCommand = command.Command{
	Name: "live",
	Parameters: []*command.Parameter{
		{
			Name:     "page",
			Required: false,
			Type:     parameters.ParameterInt,
		},
	},
	Handler: func(ctx *command.Context) command.Result {
		return ListLiveMatches(ctx, ctx.T.(*i18n.CommandLive))
	},
}

func ListLiveMatches(ctx *command.Context, t *i18n.CommandLive) command.Result {
	page := 1

	if len(ctx.Args) > 0 {
		page = ctx.Args[0].(int)
	}

	matches, err := livescore.ListLives()
	if err != nil {
		return ctx.Error(err.Error())
	}

	if len(matches) == 0 {
		return ctx.Error(t.NoMatchesLive.Str())
	}

	desc := strings.Builder{}

	totalPages := int(math.Ceil(float64(len(matches)) / float64(PageSize)))

	if page > totalPages {
		return ctx.Error(t.PageNotFound.Str())
	}

	if len(matches) > PageSize {
		desc.WriteString(t.Page.Str(page, totalPages))
		desc.WriteString("\n\n")
	}

	var pageItems []*livescore.MatchInfo
	if len(matches) > PageSize {
		pageItems = matches[(page-1)*PageSize : page*PageSize]
	} else {
		pageItems = matches
	}

	for _, match := range pageItems {
		desc.WriteString(fmt.Sprintf("%s **%s** %d x %d **%s**\n",
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
