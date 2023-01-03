package arkw

import (
	"github.com/Pauloo27/aryzona/internal/discord/model"
	dc "github.com/diamondburned/arikawa/v3/discord"
)

func buildEmbedFields(fields []*model.EmbedField) (dcgoFields []dc.EmbedField) {
	for _, field := range fields {
		dcgoFields = append(dcgoFields, dc.EmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		})
	}
	return
}

func buildEmbed(e *model.Embed) dc.Embed {
	return dc.Embed{
		Title:       e.Title,
		Description: e.Description,
		URL:         e.URL,
		Thumbnail: &dc.EmbedThumbnail{
			URL: e.ThumbnailURL,
		},
		Footer: &dc.EmbedFooter{
			Text: e.Footer,
		},
		Color: dc.Color(e.Color),
		Image: &dc.EmbedImage{
			URL: e.ImageURL,
		},
		Fields: buildEmbedFields(e.Fields),
	}
}
