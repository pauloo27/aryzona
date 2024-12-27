package disgo

import (
	dc "github.com/disgoorg/disgo/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

func buildEmbed(e *model.Embed) dc.Embed {
	builder := dc.NewEmbedBuilder().
		SetTitle(e.Title).
		SetDescription(e.Description).
		SetURL(e.URL).
		SetColor(int(e.Color)).
		SetImage(e.ImageURL).
		SetThumbnail(e.ThumbnailURL).
		SetFooter(e.Footer, "")

	for _, field := range e.Fields {
		builder.AddField(field.Name, field.Value, field.Inline)
	}

	return builder.Build()
}
