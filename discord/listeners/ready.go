package listeners

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateStreamingStatus(0, os.Getenv("DC_BOT_PRESENCE"), "https://twitch.tv/gaules")
}
