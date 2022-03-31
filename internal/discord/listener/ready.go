package listeners

import (
	"os"
	"os/user"
	"strconv"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/event"
	"github.com/Pauloo27/aryzona/internal/providers/animal"
	"github.com/Pauloo27/aryzona/internal/providers/git"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/logger"
)

func init() {
	err := discord.Bot.Listen(event.Ready, ready)
	if err != nil {
		logger.Fatal(err)
	}
}

func ready(bot discord.BotAdapter) {
	presence := os.Getenv("DC_BOT_PRESENCE")
	if presence == "" {
		presence = git.CommitMessage
	}
	err := bot.UpdatePresence(&discord.Presence{
		Title: presence,
		Type:  discord.PresenceStreaming,
		Extra: "https://twitch.tv/gaules",
	})
	if err != nil {
		logger.Error(err)
	}

	if os.Getenv("DC_BOT_ENV") == "prod" {
		c, err := discord.OpenChatWithOwner()
		if err != nil {
			logger.Fatal(err)
		}
		_, err = bot.SendEmbedMessage(c.ID(), createStartedEmbed(bot.GuildCount()))
		if err != nil {
			logger.Fatal(err)
		}
	}
}

func createStartedEmbed(guildCount int) *discord.Embed {
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

	return discord.NewEmbed().
		WithTitle(utils.Fmt("I've just started as %s@%s", userName, hostName)).
		WithColor(0xC0FFEE).
		WithImage(dogImage).
		WithField("Guild count", strconv.Itoa(guildCount)).
		WithField("Last commit", utils.Fmt("**[%s](%s/commit/%s)**", git.CommitMessage, git.RemoteRepo, git.CommitHash))
}
