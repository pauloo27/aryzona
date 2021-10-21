package voicer

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/Pauloo27/aryzona/audio/dca"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/discord/voicer/queue"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

type Voicer struct {
	usable, playing            bool
	UserID, ChannelID, GuildID *string
	Voice                      *discordgo.VoiceConnection
	Queue                      *queue.Queue
	EncodeSession              *dca.EncodeSession
	StreamingSession           *dca.StreamingSession
	lock                       *sync.Mutex
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

func (v *Voicer) registerListeners() {
	start := func(params ...interface{}) {
		if v.IsPlaying() {
			return
		}
		_ = v.Start()
	}
	v.Queue.On(queue.EventAppend, start)

	v.Queue.On(queue.EventPop, func(params ...interface{}) {
		index := params[1].(int)
		if index == 0 {
			v.EncodeSession.Cleanup()
		}
	})
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

	var chanID *string

	g, err := discord.Session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, state := range g.VoiceStates {
		if state.UserID == userID {
			chanID = &state.ChannelID
			break
		}
	}

	queue := queue.NewQueue()

	voicer = &Voicer{
		UserID: &userID, ChannelID: chanID, GuildID: &guildID, Voice: nil,
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

func (v *Voicer) IsPaused() bool {
	if v.StreamingSession == nil {
		return false
	}
	return v.StreamingSession.Paused()
}

func (v *Voicer) TogglePause() {
	if v.StreamingSession == nil {
		return
	}
	v.StreamingSession.TogglePause()
}

func (v *Voicer) Skip() {
	if v.EncodeSession == nil {
		return
	}
	v.EncodeSession.Cleanup()
}

func (v *Voicer) Connect() error {
	if !v.CanConnect() {
		return ErrCannotConnect
	}

	vc, err := discord.Session.ChannelVoiceJoin(*v.GuildID, *v.ChannelID, false, false)
	if err != nil {
		return err
	}
	v.Voice = vc
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

	v.EncodeSession.Cleanup()
	v.EncodeSession = nil

	v.Queue.Clear()
	err := v.Voice.Disconnect()
	if err != nil {
		logger.Error(err)
	}
	v.Voice = nil
	delete(voicerMapper, *(v.GuildID))
	return err
}

func (v *Voicer) IsConnected() bool {
	return v.Voice != nil
}

func (v *Voicer) IsPlaying() bool {
	return v.playing
}

func (v *Voicer) Playing() playable.Playable {
	return v.Queue.First()
}

func (v *Voicer) AppendToQueue(playable playable.Playable) error {
	v.Queue.Append(playable)
	return nil
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
		return ErrAlreadyPlaying
	}

	v.playing = true
	defer func() {
		v.playing = false
		if v.IsConnected() {
			_ = v.Disconnect()
		}
	}()

	if err := v.Connect(); err != nil {
		return err
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
		if playable == nil {
			return nil
		}

		url, err := playable.GetDirectURL()
		if err != nil {
			return err
		}

		v.EncodeSession = dca.EncodeData(url, playable.IsOppus(), playable.IsLocal())

		done = make(chan error)
		v.StreamingSession = dca.NewStream(v.EncodeSession, v.Voice, done)

		err = <-done

		v.Queue.Remove(0)

		if err == nil || err == io.EOF {
			continue
		}

		return err
	}
}
