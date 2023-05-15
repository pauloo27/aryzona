package listeners

import (
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/event"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/providers/animal"
	"github.com/pauloo27/aryzona/internal/providers/git"
	"github.com/pauloo27/logger"
)

func init() {
	err := discord.Bot.Listen(event.Ready, ready)
	if err != nil {
		logger.Fatal(err)
	}
}

func ready(bot discord.BotAdapter) {
	presence := config.Config.Presence
	if presence == "" {
		presence = git.CommitMessage
	}
	err := bot.UpdatePresence(&model.Presence{
		Title: presence,
		Type:  model.PresenceStreaming,
		Extra: "https://twitch.tv/gaules",
	})
	if err != nil {
		logger.Error(err)
	}

	if config.Config.Env == "prod" {
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

func createStartedEmbed(guildCount int) *model.Embed {
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

	return model.NewEmbed().
		WithTitle(fmt.Sprintf("I've just started as %s@%s", userName, hostName)).
		WithColor(0xC0FFEE).
		WithImage(dogImage).
		WithField("Guild count", strconv.Itoa(guildCount)).
		WithField("Last commit", fmt.Sprintf(
			"**[%s](%s/commit/%s)**",
			git.CommitMessage,
			git.RemoteRepo,
			git.CommitHash,
		))
}
