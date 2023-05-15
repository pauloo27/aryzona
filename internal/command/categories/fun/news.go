package fun

import (
	"fmt"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/news"
	"github.com/pauloo27/logger"
)

type NewsFactory func() (*news.NewsFeed, error)

var (
	Sources = map[string]NewsFactory{
		"thn":             news.GetTHNFeed,
		"cnn-world":       news.GetCNNWorldFeed,
		"cnn-tech":        news.GetCNNTechFeed,
		"cnn-top-stories": news.GetCNNTopStoriesFeed,
	}
)

func listSources() []any {
	var sources []any

	for name := range Sources {
		sources = append(sources, name)
	}

	return sources
}

var NewsCommand = command.Command{
	Name: "news",
	Parameters: []*command.CommandParameter{
		{
			Name:            "source",
			Type:            parameters.ParameterString,
			Required:        true,
			ValidValuesFunc: listSources,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandNews)

		source := ctx.Args[0].(string)
		news, err := Sources[source]()
		if err != nil {
			ctx.Error(t.SomethingWentWrong.Str())
			logger.Error(err)
			return
		}
		embed := model.NewEmbed().
			WithTitle(t.Title.Str(news.Title, news.Author)).
			WithDescription(news.Description).
			WithImage(news.ThumbnailURL)

		for _, entry := range news.Entries[:10] {
			shortDescription := entry.Description
			if len(shortDescription) > 100 {
				shortDescription = shortDescription[:97] + "..."
			}
			var postedAt string
			if entry.PostedAt != nil {
				postedAt = entry.PostedAt.Format("2006-01-02")
			}
			embed.WithField(entry.Title, fmt.Sprintf("%s | %v | %s", shortDescription, postedAt, entry.URL))
		}

		ctx.SuccessEmbed(embed)
	},
}
