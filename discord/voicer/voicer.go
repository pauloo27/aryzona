package voicer

import (
	"io"
	"sync"

	"github.com/Pauloo27/aryzona/audio/dca"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/discord/voicer/queue"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

type Voicer struct {
	playing                    bool
	UserID, ChannelID, GuildID *string
	Voice                      *discordgo.VoiceConnection
	Queue                      *queue.Queue
	EncodeSession              *dca.EncodeSession
	StreamingSession           *dca.StreamingSession
	disconnectMutex            sync.Mutex
}

var voicerMapper = map[string]*Voicer{}

func GetExistingVoicerForGuild(guildID string) *Voicer {
	return voicerMapper[guildID]
}

func (v *Voicer) registerListeners() {
	v.Queue.On(queue.EventAppend, func(params ...interface{}) {
		err := v.Start()
		if err != nil {
			logger.Error(err)
		}
	})
}

func NewVoicerForUser(userID, guildID string) (*Voicer, error) {
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

	voicer, found := voicerMapper[guildID]
	if !found {
		queue := queue.NewQueue()
		voicer = &Voicer{
			UserID: &userID, ChannelID: chanID, GuildID: &guildID, Voice: nil,
			StreamingSession: nil, EncodeSession: nil,
			disconnectMutex: sync.Mutex{},
			Queue:           queue,
		}
		voicer.registerListeners()
		voicerMapper[guildID] = voicer
	}
	return voicer, nil

}

func (v *Voicer) CanConnect() bool {
	return v.ChannelID != nil
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
	delete(voicerMapper, *(v.GuildID))
	v.disconnectMutex.Lock()
	defer v.disconnectMutex.Unlock()
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

func (v *Voicer) Start() error {
	if v.IsPlaying() {
		return ErrAlreadyPlaying
	}

	v.playing = true
	defer func() {
		v.playing = false
	}()

	if err := v.Connect(); err != nil {
		return err
	}

	// play a simple "pre connect" sound
	v.EncodeSession = dca.EncodeData("./assets/radio_start.wav", false, true)

	for {
		playable := v.Queue.First()
		if playable == nil {
			return nil
		}

		url, err := playable.GetDirectURL()
		if err != nil {
			return err
		}

		if err := v.Voice.Speaking(true); err != nil {
			return err
		}

		done := make(chan error)
		v.StreamingSession = dca.NewStream(v.EncodeSession, v.Voice, done)

		err = <-done
		if err != nil && err != io.EOF {
			logger.Error(err)
		}

		v.EncodeSession = dca.EncodeData(url, playable.IsOppus(), playable.IsLocal())

		done = make(chan error)
		v.StreamingSession = dca.NewStream(v.EncodeSession, v.Voice, done)

		err = <-done
		if v.IsConnected() {
			disconnectErr := v.Disconnect()
			if disconnectErr != nil {
				return utils.Wrap(disconnectErr.Error(), err)
			}
		}

		if err == io.EOF {
			continue
		}
		return err
	}
}
