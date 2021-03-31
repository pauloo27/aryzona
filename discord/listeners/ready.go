package listeners

import (
	"os"

	"github.com/Pauloo27/aryzona/utils"
	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateStreamingStatus(0, utils.Fmt("last commit: %s", os.Getenv("DC_BOT_PRESENCE")), "https://twitch.tv/gaules")
}
