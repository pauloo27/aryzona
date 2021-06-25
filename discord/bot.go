package discord

import (
	"time"

	"github.com/Pauloo27/aryzona/discord/listeners"
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

func Disconnect() {
	Session.Close()
}

func AddDefaultListeners() {
	Session.AddHandler(listeners.MessageCreate)
	Session.AddHandler(listeners.Ready)
}
