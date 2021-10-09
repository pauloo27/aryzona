package listeners

import (
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/bwmarrin/discordgo"
)

func init() {
	discord.Listen(VoiceChannelDisconnect)
}

func countUsersInChannel(guildID, channelID string) (count int) {
	g, err := discord.Session.State.Guild(guildID)
	if err != nil {
		return 0
	}
	for _, voice := range g.VoiceStates {
		if voice.ChannelID == channelID {
			count++
		}
	}
	return
}

func VoiceChannelDisconnect(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	if s.State.User.ID == e.UserID {
		return
	}
	v := voicer.GetExistingVoicerForGuild(e.GuildID)
	if v == nil || v.ChannelID == nil {
		return
	}
	if e.BeforeUpdate.ChannelID != *v.ChannelID {
		return
	}
	userCount := countUsersInChannel(e.GuildID, e.ChannelID)
	if userCount > 1 {
		return
	}
	// TODO: schedule disconnect - Instead of disconnecting right away, wait around
	// 30 seconds because maybe the user will come back. If they dont, then we say
	// bye bye =)

	// TODO: how do "unschedule" when the user come back? Btw, handle that in the current
	// listener
}
