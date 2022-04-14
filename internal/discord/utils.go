package discord

import (
	"fmt"
	"os"
)

func AsMention(userID string) string {
	return fmt.Sprintf("<@%s>", userID)
}

func OpenChatWithOwner() (Channel, error) {
	return Bot.OpenChannelWithUser(os.Getenv("DC_BOT_OWNER_ID"))
}
