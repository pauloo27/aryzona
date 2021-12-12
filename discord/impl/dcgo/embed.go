package dcgo

import (
	"github.com/Pauloo27/aryzona/discord"
	"github.com/bwmarrin/discordgo"
)

func buildEmbedFields(fields []*discord.EmbedField) (dcgoFields []*discordgo.MessageEmbedField) {
	for _, field := range fields {
		dcgoFields = append(dcgoFields, &discordgo.MessageEmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		})
	}
	return
}

func buildEmbed(e *discord.Embed) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       e.Title,
		Description: e.Description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: e.ThumbnailURL,
		},
		Color: e.Color,
		Image: &discordgo.MessageEmbedImage{
			URL: e.ImageURL,
		},
		Fields: buildEmbedFields(e.Fields),
	}
}
