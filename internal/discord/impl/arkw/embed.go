package arkw

import (
	"github.com/pauloo27/aryzona/internal/discord/model"
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
	var image *dc.EmbedImage
	var thumbnail *dc.EmbedThumbnail
	var footer *dc.EmbedFooter

	if e.ImageURL != "" {
		image = &dc.EmbedImage{
			URL: e.ImageURL,
		}
	}

	if e.ThumbnailURL != "" {
		thumbnail = &dc.EmbedThumbnail{
			URL: e.ThumbnailURL,
		}
	}

	if e.Footer != "" {
		footer = &dc.EmbedFooter{
			Text: e.Footer,
		}
	}

	return dc.Embed{
		Title:       e.Title,
		Description: e.Description,
		URL:         e.URL,
		Thumbnail:   thumbnail,
		Footer:      footer,
		Color:       dc.Color(e.Color),
		Image:       image,
		Fields:      buildEmbedFields(e.Fields),
	}
}
