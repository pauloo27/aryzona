package utils

import (
	"github.com/bwmarrin/discordgo"
)

type EmbedBuilder struct {
	embed *discordgo.MessageEmbed
}

func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{embed: &discordgo.MessageEmbed{}}
}

func (e *EmbedBuilder) Color(color int) *EmbedBuilder {
	e.embed.Color = color
	return e
}

func (e *EmbedBuilder) Title(title string) *EmbedBuilder {
	e.embed.Title = title
	return e
}

func (e *EmbedBuilder) Description(description string) *EmbedBuilder {
	e.embed.Description = description
	return e
}

func (e *EmbedBuilder) Field(name, value string) *EmbedBuilder {
	e.embed.Fields = append(e.embed.Fields, &discordgo.MessageEmbedField{
		Name: name,
		Value: value,
	})
	return e
}

func (e *EmbedBuilder) Build() *discordgo.MessageEmbed {
	return e.embed
}
