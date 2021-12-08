package discord

type Message struct {
	ID      string
	Author  *User
	Channel *Channel
	Content string
}
