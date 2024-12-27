package disgo

import (
	"github.com/pauloo27/aryzona/internal/discord/model"
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
