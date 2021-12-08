package discord

type Message interface {
	ID() string
	Author() User
	Channel() Channel
	Content() string
}
