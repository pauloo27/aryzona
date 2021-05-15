package listeners

import (
	"os"
	"os/user"

	"github.com/Pauloo27/aryzona/git"
	"github.com/Pauloo27/aryzona/provider/animal"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/bwmarrin/discordgo"
)

func createStartedEmbed() *discordgo.MessageEmbed {
	dogImage, err := animal.GetRandomDogImage()
	if err != nil {
		dogImage = "https://http.cat/500"
	}

	userName := "??"
	user, err := user.Current()
	if err == nil {
		userName = user.Username
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "??"
	}

	return &discordgo.MessageEmbed{
		Title: utils.Fmt("I've just started in %s@%s", userName, hostname),
		Color: 0xC0FFEE,
		Image: &discordgo.MessageEmbedImage{URL: dogImage},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Last commit", Value: utils.Fmt(
				"**[%s](%s/commit/%s)**", git.CommitMessage, git.RemoteRepo, git.CommitHash,
			)},
		},
	}
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	presence := os.Getenv("DC_BOT_PRESENCE")
	if presence == "" {
		presence = git.CommitMessage
	}
	s.UpdateStreamingStatus(0, presence, "https://twitch.tv/gaules")

	c, err := s.UserChannelCreate(os.Getenv("DC_BOT_OWNER_ID"))
	utils.HandleFatal(err)

	_, err = s.ChannelMessageSendEmbed(c.ID, createStartedEmbed())
	utils.HandleFatal(err)
}
