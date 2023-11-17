package listener

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/core/scheduler"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/event"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
)

func init() {
	err := discord.Bot.Listen(event.VoiceStateUpdated, voiceUpdate)
	if err != nil {
		panic(err)
	}
}

func voiceUpdate(bot discord.BotAdapter, user model.User, prevCh, curCh model.VoiceChannel) {
	if curCh != nil {
		onConnect(bot, user, prevCh, curCh)
	}

	if prevCh != nil {
		onDisconnect(bot, user, prevCh, curCh)
	}
}

func onConnect(bot discord.BotAdapter, user model.User, prevCh, curCh model.VoiceChannel) {
	if isUserBot(user) {
		onBotConnected(prevCh, curCh)
		return
	}
	checkConnectedChannel(bot, curCh)
}

func onDisconnect(bot discord.BotAdapter, user model.User, prevCh, curCh model.VoiceChannel) {
	if isUserBot(user) {
		onBotDisconnected(prevCh, curCh)
		return
	}
	checkDisconnectedChannel(bot, prevCh)
}

func checkConnectedChannel(bot discord.BotAdapter, ch model.VoiceChannel) {
	v := voicer.GetExistingVoicerForGuild(ch.Guild().ID())

	if !isVoicerValid(v) || *v.ChannelID != ch.ID() {
		return
	}

	if bot.CountUsersInVoiceChannel(ch) > 1 {
		unsheduleDisconnect(v)
	}
}

func checkDisconnectedChannel(bot discord.BotAdapter, prevCh model.VoiceChannel) {
	v := voicer.GetExistingVoicerForGuild(prevCh.Guild().ID())
	if !isVoicerValid(v) || *v.ChannelID != prevCh.ID() {
		return
	}

	if bot.CountUsersInVoiceChannel(prevCh) <= 1 {
		scheduleDisconnect(v)
	}
}

func onBotConnected(prevCh, curCh model.VoiceChannel) {
	v := voicer.GetExistingVoicerForGuild(curCh.Guild().ID())
	if !isVoicerValid(v) {
		return
	}

	// if the bot connected to the channel it's expected to, there's nothing to do
	if *v.ChannelID == curCh.ID() {
		return
	}

	*v.ChannelID = curCh.ID()

	usersInCh := discord.Bot.CountUsersInVoiceChannel(curCh)

	if usersInCh > 1 {
		unsheduleDisconnect(v)
	} else {
		scheduleDisconnect(v)
	}
}

func onBotDisconnected(prevCh, curCh model.VoiceChannel) {
	v := voicer.GetExistingVoicerForGuild(prevCh.Guild().ID())
	if isVoicerValid(v) && curCh == nil {
		_ = v.Disconnect()
	}
}

func isUserBot(user model.User) bool {
	self, err := discord.Bot.Self()
	if err != nil {
		return false
	}

	return self.ID() == user.ID()
}

func isVoicerValid(v *voicer.Voicer) bool {
	return v != nil && v.ChannelID != nil && v.GuildID != nil
}

func scheduleDisconnect(v *voicer.Voicer) {
	task := scheduler.NewRunLaterTask(
		60*time.Second,
		func(params ...any) {
			if err := v.Disconnect(); err != nil {
				slog.Error("Cannot disconnect empty channel", tint.Err(err))
			}
		},
	)

	scheduler.Schedule(fmt.Sprintf("voice_disconnect_%s", *v.GuildID), task)
}

func unsheduleDisconnect(v *voicer.Voicer) {
	scheduler.Unschedule(fmt.Sprintf("voice_disconnect_%s", *v.GuildID))
}
