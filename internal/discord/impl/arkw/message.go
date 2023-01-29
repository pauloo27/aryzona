package arkw

import (
	"github.com/Pauloo27/aryzona/internal/discord/model"
	dc "github.com/diamondburned/arikawa/v3/discord"
)

type Message struct {
	ch      Channel
	id      string
	author  User
	content string
}

func (m Message) ID() string {
	return m.id
}

func (m Message) Content() string {
	return m.content
}

func (m Message) Channel() model.TextChannel {
	return m.ch
}

func (m Message) Author() model.User {
	return m.author
}

func buildMessage(id string, ch Channel, author User, content string) Message {
	return Message{
		id:      id,
		ch:      ch,
		author:  author,
		content: content,
	}
}

type ArkwComplexMessage struct {
	embed      []dc.Embed
	components dc.ContainerComponents
}

func prepareComplexMessage(message *model.ComplexMessage) *ArkwComplexMessage {
	embeds := make([]dc.Embed, len(message.Embeds))
	if len(message.Embeds) > 0 {
		for i, embed := range message.Embeds {
			embeds[i] = buildEmbed(embed)
		}
	}

	var componentsPtr dc.ContainerComponents

	if len(message.Components) > 0 {
		components := buildComponents(message.Components)
		row := dc.ActionRowComponent(components)
		componentsPtr = dc.ContainerComponents{&row}
	}

	return &ArkwComplexMessage{
		components: componentsPtr,
		embed:      embeds,
	}
}
