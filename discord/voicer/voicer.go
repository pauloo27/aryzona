package voicer

import (
	"errors"

	"github.com/Pauloo27/aryzona/audio"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/bwmarrin/discordgo"
)

type Voicer struct {
	ChannelID, GuildID *string
	Voice              *discordgo.VoiceConnection
	Playing            *audio.Playable
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
	return &Voicer{chanID, &guildID, nil, nil}, nil
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
	}

	return nil
}
