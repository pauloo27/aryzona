package discord

import (
	"github.com/Pauloo27/aryzona/discord/listeners"
	"github.com/bwmarrin/discordgo"
)

var Discord *discordgo.Session

func Create(token string) error {
	var err error
	Discord, err = discordgo.New("Bot " + token)
	return err
}

func Connect() error {
	return Discord.Open()
}

func Disconnect() {
	Discord.Close()
}

func AddDefaultListeners() {
	Discord.AddHandler(listeners.MessageCreate)
}
