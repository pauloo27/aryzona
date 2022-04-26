package model

type Message interface {
	ID() string
	Author() User
	Channel() TextChannel
	Content() string
}
