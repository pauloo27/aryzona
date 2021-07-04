package listeners

import (
	"os"
	"os/user"
	"strconv"

	"github.com/Pauloo27/aryzona/git"
	"github.com/Pauloo27/aryzona/providers/animal"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/bwmarrin/discordgo"
)

func createStartedEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
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
			{
				Name:  "Guild count",
				Value: strconv.Itoa(len(s.State.Guilds)),
			},
			{
				Name: "Last commit", Value: utils.Fmt(
					"**[%s](%s/commit/%s)**", git.CommitMessage, git.RemoteRepo, git.CommitHash,
				),
			},
		},
	}
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	presence := os.Getenv("DC_BOT_PRESENCE")
	if presence == "" {
		presence = utils.Fmt("Last commit: %s", git.CommitMessage)
	}
	s.UpdateStreamingStatus(0, presence, "https://twitch.tv/gaules")

	if os.Getenv("DC_BOT_ENV") == "prod" {
		c, err := s.UserChannelCreate(os.Getenv("DC_BOT_OWNER_ID"))
		utils.HandleFatal(err)
		_, err = s.ChannelMessageSendEmbed(c.ID, createStartedEmbed(s))
		utils.HandleFatal(err)
	}
}
