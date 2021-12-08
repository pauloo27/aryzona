package discord

import (
	"os"

	"github.com/Pauloo27/aryzona/utils"
)

func AsMention(userID string) string {
	return utils.Fmt("<@%s>", userID)
}

func OpenChatWithOwner() (*Channel, error) {
	return Bot.OpenChannelWithUser(os.Getenv("DC_BOT_OWNER_ID"))
}
