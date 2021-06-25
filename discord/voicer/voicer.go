package voicer

import (
	"errors"
	"time"

	"github.com/Pauloo27/aryzona/audio"
	"github.com/Pauloo27/aryzona/audio/dca"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/logger"
	"github.com/bwmarrin/discordgo"
)

type Voicer struct {
	ChannelID, GuildID *string
	Voice              *discordgo.VoiceConnection
	Playing            *audio.Playable
	EncodeSession      *dca.EncodeSession
	StreamingSession   *dca.StreamingSession
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
	return &Voicer{chanID, &guildID, nil, nil, nil, nil}, nil
}

func (v *Voicer) CanConnect() bool {
	return v.ChannelID != nil
}

func (v *Voicer) Connect() error {
	if !v.CanConnect() {
		return errors.New("Cannot connect")
	}
	vc, err := discord.Session.ChannelVoiceJoin(*v.GuildID, *v.ChannelID, false, false)
	if err != nil {
		return err
	}
	v.Voice = vc
	return nil
}

func (v *Voicer) Disconnect() error {
	err := v.Voice.Disconnect()
	v.Voice = nil
	return err
}

func (v *Voicer) IsConnected() bool {
	return v.Voice != nil
}

func (v *Voicer) Play(playable audio.Playable) error {
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

	return <-done
}
