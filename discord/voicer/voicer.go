package voicer

import (
	"sync"
	"time"

	"github.com/Pauloo27/aryzona/audio"
	"github.com/Pauloo27/aryzona/audio/dca"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/logger"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/bwmarrin/discordgo"
)

type Voicer struct {
	ChannelID, GuildID *string
	Voice              *discordgo.VoiceConnection
	Playing            *audio.Playable
	EncodeSession      *dca.EncodeSession
	StreamingSession   *dca.StreamingSession
	disconnectMutex    sync.Mutex
}

var voicerMapper = map[string]*Voicer{}

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
		voicer = &Voicer{chanID, &guildID, nil, nil, nil, nil, sync.Mutex{}}
		voicerMapper[guildID] = voicer
	}
	return voicer, nil

}

func (v *Voicer) CanConnect() bool {
	return v.ChannelID != nil
}

func (v *Voicer) Connect() error {
	if !v.CanConnect() {
		return ERR_CANNOT_CONNECT
	}

	vc, err := discord.Session.ChannelVoiceJoin(*v.GuildID, *v.ChannelID, false, false)
	if err != nil {
		return err
	}
	v.Voice = vc
	return nil
}

func (v *Voicer) Disconnect() error {
	v.disconnectMutex.Lock()
	defer v.disconnectMutex.Unlock()
	if !v.IsConnected() {
		return nil
	}

	v.StreamingSession = nil

	v.EncodeSession.Cleanup()
	v.EncodeSession.Stop()
	v.EncodeSession = nil

	v.Playing = nil
	err := v.Voice.Disconnect()
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
		return ERR_ALREADY_PLAYING
	}
	if !v.IsConnected() {
		if err := v.Connect(); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	v.Playing = &playable
	url, err := playable.GetDirectURL()
	if err != nil {
		return err
	}
	if err := v.Voice.Speaking(true); err != nil {
		return err
	}
	logger.Debugf("playing %s", url)

	v.EncodeSession = dca.EncodeData(url, playable.IsOppus())
	defer v.EncodeSession.Cleanup()

	done := make(chan error)
	v.StreamingSession = dca.NewStream(v.EncodeSession, v.Voice, done)

	err = <-done
	disconnectErr := v.Disconnect()
	if disconnectErr != nil {
		return utils.Wrap(disconnectErr.Error(), err)
	}

	return err
}
