package discord

import (
	"time"

	"github.com/Pauloo27/aryzona/discord/listeners"
	"github.com/bwmarrin/discordgo"
)

var Discord *discordgo.Session
var StartedAt time.Time

func Create(token string) error {
	var err error
	Discord, err = discordgo.New("Bot " + token)
	return err
}

func Connect() error {
	StartedAt = time.Now()
	return Discord.Open()
}

func Disconnect() {
	Discord.Close()
}

func AddDefaultListeners() {
	Discord.AddHandler(listeners.MessageCreate)
	Discord.AddHandler(listeners.Ready)
}
