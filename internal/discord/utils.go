package discord

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/discord/model"
)

func AsMention(userID string) string {
	return fmt.Sprintf("<@%s>", userID)
}

func OpenChatWithOwner() (model.TextChannel, error) {
	return Bot.OpenChannelWithUser(config.Config.OwnerID)
}
