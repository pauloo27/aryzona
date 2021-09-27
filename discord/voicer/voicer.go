package voicer

import (
	"io"
	"sync"

	"github.com/Pauloo27/aryzona/audio"
	"github.com/Pauloo27/aryzona/audio/dca"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

type Voicer struct {
	UserID, ChannelID, GuildID *string
	Voice                      *discordgo.VoiceConnection
	Playing                    *audio.Playable
	EncodeSession              *dca.EncodeSession
	StreamingSession           *dca.StreamingSession
	disconnectMutex            sync.Mutex
}

var voicerMapper = map[string]*Voicer{}

func GetExistingVoicerForGuild(guildID string) *Voicer {
	return voicerMapper[guildID]
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
		voicer = &Voicer{&userID, chanID, &guildID, nil, nil, nil, nil, sync.Mutex{}}
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

	v.Playing = nil
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
	return v.Playing != nil
}

func (v *Voicer) Play(playable audio.Playable) error {
	if v.IsPlaying() {
		return ErrAlreadyPlaying
	}

	if err := v.Connect(); err != nil {
		return err
	}

	v.Playing = &playable

	url, err := playable.GetDirectURL()
	if err != nil {
		return err
	}

	if err := v.Voice.Speaking(true); err != nil {
		return err
	}

	// play a simple "pre connect" sound
	v.EncodeSession = dca.EncodeData("./assets/radio_start.wav", false, true)

	done := make(chan error)
	v.StreamingSession = dca.NewStream(v.EncodeSession, v.Voice, done)

	err = <-done
	if err != nil {
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
		return nil
	}
	return err
}
