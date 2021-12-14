package arkw

import (
	"github.com/Pauloo27/aryzona/internal/discord"
	dc "github.com/diamondburned/arikawa/v3/discord"
)

func buildEmbedFields(fields []*discord.EmbedField) (dcgoFields []dc.EmbedField) {
	for _, field := range fields {
		dcgoFields = append(dcgoFields, dc.EmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		})
	}
	return
}

func buildEmbed(e *discord.Embed) dc.Embed {
	return dc.Embed{
		Title:       e.Title,
		Description: e.Description,
		Thumbnail: &dc.EmbedThumbnail{
			URL: e.ThumbnailURL,
		},
		Color: dc.Color(e.Color),
		Image: &dc.EmbedImage{
			URL: e.ImageURL,
		},
		Fields: buildEmbedFields(e.Fields),
	}
}
