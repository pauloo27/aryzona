package model

type Message interface {
	ID() string
	Author() User
	Channel() TextChannel
	Content() string
}

type ComplexMessage struct {
	Content    string
	Embeds     []*Embed
	Components []MessageComponent
	ReplyTo    Message
}
