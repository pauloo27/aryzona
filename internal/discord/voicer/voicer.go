package voicer

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/Pauloo27/aryzona/internal/audio/dca"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/queue"
	"github.com/Pauloo27/logger"
)

type Voicer struct {
	Voice                      discord.VoiceConnection
	Queue                      *queue.Queue
	EncodeSession              *dca.EncodeSession
	UserID, ChannelID, GuildID *string
	StreamingSession           *dca.StreamingSession
	lock                       *sync.Mutex
	usable, playing            bool
}

var (
	voicerMapper = map[string]*Voicer{}
)

func GetExistingVoicerForGuild(guildID string) *Voicer {
	return voicerMapper[guildID]
}

func (v *Voicer) Lock() {
	v.lock.Lock()
}

func (v *Voicer) Unlock() {
	v.lock.Unlock()
}

func NewVoicerForUser(userID, guildID string) (*Voicer, error) {
	voicer, found := voicerMapper[guildID]

	if found {
		voicer.Lock()
		defer voicer.Unlock()
		// not usable means fully disconnected or going to disconnect
		// why? cuz doing ,play then ,stop ,play very quickly broke the bot =(
		if voicer.usable {
			return voicer, nil
		}
	}

	g, err := discord.Bot.OpenGuild(guildID)
	if err != nil {
		return nil, err
	}

	vc, err := discord.Bot.FindUserVoiceState(g.ID(), userID)
	if err != nil {
		return nil, err
	}
	chanID := vc.Channel().ID()

	queue := queue.NewQueue()

	voicer = &Voicer{
		UserID: &userID, ChannelID: &chanID, GuildID: &guildID, Voice: nil,
		StreamingSession: nil, EncodeSession: nil,
		lock:  &sync.Mutex{},
		Queue: queue, usable: true,
	}

	voicer.registerListeners()
	voicerMapper[guildID] = voicer

	return voicer, nil

}

func (v *Voicer) CanConnect() bool {
	return v.ChannelID != nil
}

func (v *Voicer) Connect() error {
	if !v.CanConnect() {
		return errors.New("cannot connect")
	}

	if v.IsConnected() {
		return nil
	}

	vc, err := discord.Bot.JoinVoiceChannel(*v.GuildID, *v.ChannelID)
	if err != nil {
		return err
	}
	v.Voice = vc

	// avoid connecting but never playing a thing
	v.scheduleEmptyQueue()

	return nil
}

func (v *Voicer) Disconnect() error {
	v.Lock()
	defer v.Unlock()

	// usable = fully disconnected or going to disconnect
	v.usable = false

	if !v.IsConnected() {
		return nil
	}

	v.StreamingSession = nil

	if v.EncodeSession != nil {
		v.EncodeSession.Cleanup()
		v.EncodeSession = nil
	}

	var err error
	if v.Queue != nil {
		v.Queue.Clear()
		err = v.Voice.Disconnect()
		if err != nil {
			logger.Error(err)
		}
		v.Voice = nil
	}

	delete(voicerMapper, *(v.GuildID))
	return err
}

func (v *Voicer) IsConnected() bool {
	return v.Voice != nil
}

func (v *Voicer) IsPlaying() bool {
	return v.playing
}

func (v *Voicer) GetPosition() (time.Duration, error) {
	if v == nil || v.Queue.First() == nil {
		// TODO: create an errore?
		return 0, errors.New("nothing playing")
	}
	return time.Duration(v.StreamingSession.PlaybackPosition()), nil
}

func (v *Voicer) Start() error {
	if v.IsPlaying() {
		return errors.New("already playing")
	}

	v.playing = true
	defer func() {
		if v.Voice != nil && v.usable {
			_ = v.Voice.Speaking(false)
		}
		v.playing = false
		v.scheduleEmptyQueue()
	}()

	if !v.IsConnected() {
		if err := v.Connect(); err != nil {
			return err
		}
	}

	if err := v.Voice.Speaking(true); err != nil {
		return err
	}

	// play a simple "pre connect" sound
	v.EncodeSession = dca.EncodeData("./assets/radio_start.wav", false, true)
	done := make(chan error)
	v.StreamingSession = dca.NewStream(v.EncodeSession, v.Voice, done)

	err := <-done
	if err != nil && err != io.EOF {
		logger.Error(err)
	}

	for {
		playable := v.Queue.First()
		if playable == nil || !v.usable {
			return nil
		}

		url, err := playable.GetDirectURL()
		if err != nil {
			return err
		}

		v.EncodeSession = dca.EncodeData(url, playable.IsOpus(), playable.IsLocal())

		if err := v.Voice.Speaking(true); err != nil {
			return err
		}

		done := make(chan error)
		v.StreamingSession = dca.NewStream(v.EncodeSession, v.Voice, done)

		err = <-done

		v.Queue.Remove(0)

		if err == nil || err == io.EOF {
			continue
		}

		return err
	}
}
