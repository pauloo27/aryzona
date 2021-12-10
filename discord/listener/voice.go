package listeners

import (
	"time"

	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/event"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/aryzona/utils/scheduler"
	"github.com/Pauloo27/logger"
)

func init() {
	err := discord.Bot.Listen(event.VoiceStateUpdated, voiceUpdate)
	if err != nil {
		panic(err)
	}
}

func voiceUpdate(bot discord.BotAdapter, user discord.User, prevCh, curCh discord.VoiceChannel) {
	self, err := bot.Self()
	if err != nil {
		return
	}

	if self.ID() == user.ID() {
		return
	}
	var v *voicer.Voicer
	if curCh != nil {
		v = voicer.GetExistingVoicerForGuild(curCh.Guild().ID())
	}

	if v == nil || v.ChannelID == nil {
		return
	}
	voicerChan := *v.ChannelID

	if (prevCh != nil && prevCh.ID() == voicerChan) && (curCh != nil && curCh.ID() != prevCh.ID()) {
		onDisconnect(bot, curCh, v)
		return
	}

	if curCh.ID() == voicerChan {
		onConnect(bot, curCh)
		return
	}
}

func onConnect(bot discord.BotAdapter, ch discord.VoiceChannel) {
	if bot.CountUsersInVoiceChannel(ch) <= 1 {
		return
	}

	scheduler.Unschedule(utils.Fmt("voice_disconnect_%s", ch.Guild().ID()))
}

func onDisconnect(bot discord.BotAdapter, ch discord.VoiceChannel, v *voicer.Voicer) {
	if bot.CountUsersInVoiceChannel(ch) > 1 {
		return
	}

	task := scheduler.NewRunLaterTask(
		30*time.Second,
		func(params ...interface{}) {
			if err := v.Disconnect(); err != nil {
				logger.Errorf("cannot disconnect empty channel: %v", err)
			}
		},
	)

	scheduler.Schedule(utils.Fmt("voice_disconnect_%s", ch.Guild().ID()), task)
}
