package discord

import (
	"fmt"
	"os"

	"github.com/Pauloo27/aryzona/internal/discord/model"
)

func AsMention(userID string) string {
	return fmt.Sprintf("<@%s>", userID)
}

func OpenChatWithOwner() (model.TextChannel, error) {
	return Bot.OpenChannelWithUser(os.Getenv("DC_BOT_OWNER_ID"))
}
