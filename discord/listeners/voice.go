package listeners

import (
	"time"

	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/aryzona/utils/scheduler"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

func init() {
	discord.Listen(VoiceUpdate)
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

func onConnect(s *discordgo.Session, e *discordgo.VoiceStateUpdate, v *voicer.Voicer) {
	userCount := countUsersInChannel(e.GuildID, e.ChannelID)
	if userCount <= 1 {
		return
	}

	scheduler.Unschedule(utils.Fmt("voice_disconnect_%s", e.GuildID))
}

func onDisconnect(s *discordgo.Session, e *discordgo.VoiceStateUpdate, v *voicer.Voicer) {
	userCount := countUsersInChannel(e.GuildID, e.ChannelID)
	if userCount > 1 {
		return
	}

	task := scheduler.Task{
		Time: time.Now().Add(30 * time.Second),
		Callback: func(parmas ...interface{}) {
			err := v.Disconnect()
			if err != nil {
				logger.Errorf("cannot disconnect empty channel: %v", err)
			}
		},
	}
	scheduler.Schedule(utils.Fmt("voice_disconnect_%s", e.GuildID), &task)
}

func VoiceUpdate(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	if s.State.User.ID == e.UserID {
		return
	}
	v := voicer.GetExistingVoicerForGuild(e.GuildID)

	if v == nil || v.ChannelID == nil {
		return
	}
	voicerChan := *v.ChannelID

	prevChan := e.BeforeUpdate.ChannelID
	currentChan := e.VoiceState.ChannelID

	if prevChan == voicerChan && currentChan != prevChan {
		onDisconnect(s, e, v)
		return
	}

	if currentChan == voicerChan {
		onConnect(s, e, v)
		return
	}
}
