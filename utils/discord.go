package utils

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func AsMention(userID string) string {
	return Fmt("<@%s>", userID)
}

func OpenChatWithOwner(s *discordgo.Session) (*discordgo.Channel, error) {
	return s.UserChannelCreate(os.Getenv("DC_BOT_OWNER_ID"))
}
