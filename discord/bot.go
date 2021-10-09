package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

var Session *discordgo.Session
var StartedAt time.Time

func Create(token string) error {
	var err error
	Session, err = discordgo.New("Bot " + token)
	return err
}

func Connect() error {
	StartedAt = time.Now()
	return Session.Open()
}

func Disconnect() error {
	return Session.Close()
}

var listeners []interface{}

func Listen(listener interface{}) {
	listeners = append(listeners, listener)
}

func RegisterListeners() {
	for _, listener := range listeners {
		Session.AddHandler(listener)
	}
}
