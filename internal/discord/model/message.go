package model

type Message interface {
	ID() string
	Author() User
	Channel() TextChannel
	Content() string
}

type MessageComponentRow struct {
	Components []MessageComponent
}

type ComplexMessage struct {
	Content       string
	Embeds        []*Embed
	ComponentRows []MessageComponentRow
	ReplyTo       Message
}
