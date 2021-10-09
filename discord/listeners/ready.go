package listeners

import (
	"os"
	"os/user"
	"strconv"

	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/git"
	"github.com/Pauloo27/aryzona/providers/animal"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/aryzona/utils/errore"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

func init() {
	discord.Listen(Ready)
}

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

	hostName, err := os.Hostname()
	if err != nil {
		hostName = "??"
	}

	return &discordgo.MessageEmbed{
		Title: utils.Fmt("I've just started as %s@%s", userName, hostName),
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
		presence = utils.Fmt("%s", git.CommitMessage)
	}
	err := s.UpdateStreamingStatus(0, presence, "https://twitch.tv/gaules")
	if err != nil {
		logger.Error(err)
	}

	if os.Getenv("DC_BOT_ENV") == "prod" {
		c, err := s.UserChannelCreate(os.Getenv("DC_BOT_OWNER_ID"))
		errore.HandleFatal(err)
		_, err = s.ChannelMessageSendEmbed(c.ID, createStartedEmbed(s))
		errore.HandleFatal(err)
	}
}
