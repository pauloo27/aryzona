package listeners

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	presence := os.Getenv("DC_BOT_PRESENCE")
	if presence == "" {
		presence = ",help"
	}
	s.UpdateStreamingStatus(0, presence, "https://twitch.tv/gaules")
}
