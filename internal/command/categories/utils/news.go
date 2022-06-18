package utils

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/providers/news"
)

type NewsFactory func() (*news.NewsFeed, error)

var (
	Sources = map[string]NewsFactory{
		"thn": news.GetTHNFeed,
	}
)

func listSources() []interface{} {
	var sources []interface{}

	for name := range Sources {
		sources = append(sources, name)
	}

	return sources
}

var NewsCommand = command.Command{
	Name: "news", Aliases: []string{"noticias", "notícias"},
	Description: "Get recent news",
	Parameters: []*command.CommandParameter{
		{
			Name:            "source",
			Description:     "Source Name",
			Type:            parameters.ParameterString,
			Required:        true,
			ValidValuesFunc: listSources,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		source := ctx.Args[0].(string)
		news, err := Sources[source]()
		if err != nil {
			ctx.Error(fmt.Sprintf("Error getting news: %s", err))
			return
		}
		embed := discord.NewEmbed().
			WithTitle(fmt.Sprintf("News - %s by %s", news.Title, news.Author)).
			WithDescription(news.Description).
			WithImage(news.ThumbnailURL)

		for _, entry := range news.Entries[:10] {
			shortDescription := entry.Description
			if len(shortDescription) > 200 {
				shortDescription = shortDescription[:197] + "..."
			}
			postedAt := entry.PostedAt.Format("02/01/2006")
			embed.WithField(entry.Title, fmt.Sprintf("%s | %v | %s", shortDescription, postedAt, entry.URL))
		}

		ctx.SuccessEmbed(embed)
	},
}
