package listener

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"strconv"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/event"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/providers/animal"
	"github.com/pauloo27/aryzona/internal/providers/git"
)

func init() {
	err := discord.Bot.Listen(event.Ready, ready)
	if err != nil {
		slog.Error("Cannot listen to ready event", tint.Err(err))
		os.Exit(1)
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
		slog.Error("Cannot update bot presence", tint.Err(err))
	}

	if config.Config.Env != "local" {
		c, err := discord.OpenChatWithOwner()
		if err != nil {
			slog.Error("Cannot open owner chat", tint.Err(err))
			os.Exit(1)
		}
		_, err = bot.SendEmbedMessage(c.ID(), createStartedEmbed(bot.GuildCount()))
		if err != nil {
			slog.Error("Cannot send start embed to bot owner", tint.Err(err))
			os.Exit(1)
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
