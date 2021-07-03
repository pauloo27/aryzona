package radio

import (
	"strings"
)

type RadioType struct {
	Name         string
	IsOppus      bool
	GetDirectURL func(url string) string
}

type RadioChannel struct {
	Id, Name, URL string
	Type          *RadioType
}

func (c RadioChannel) CanPause() bool {
	return true
}

func (c RadioChannel) Pause() error {
	return nil
}

func (c RadioChannel) Unpause() error {
	return nil
}

func (c RadioChannel) TogglePause() error {
	return nil
}

func (c RadioChannel) GetDirectURL() (string, error) {
	return c.Type.GetDirectURL(c.URL), nil
}

func (c RadioChannel) IsOppus() bool {
	return c.Type.IsOppus
}

var radios = []*RadioChannel{
	{
		Id:   "hunter-pop",
		Name: "Rádio Hunter POP",
		URL:  "https://hls.hunter.fm/pop/192.m3u8",
		Type: M3uRadio,
	},
	{
		Id:   "hunter-pisadinha",
		Name: "Rádio Hunter Pisadinha",
		URL:  "https://hls.hunter.fm/pisadinha/320.m3u8",
		Type: M3uRadio,
	},
	{
		Id:   "hunter-tropical",
		Name: "Rádio Hunter Tropical",
		URL:  "https://hls.hunter.fm/tropical/192.m3u8",
		Type: M3uRadio,
	},
	{
		Id:   "hunter-80s",
		Name: "Rádio Hunter Anos 80",
		URL:  "https://hls.hunter.fm/80s/192.m3u8",
		Type: M3uRadio,
	},
	{
		Id:   "hunter-rock",
		Name: "Rádio Hunter Rock",
		URL:  "https://hls.hunter.fm/rock/192.m3u8",
		Type: M3uRadio,
	},
	{
		Id:   "hunter-lofi",
		Name: "Rádio Hunter LoFi",
		URL:  "https://hls.hunter.fm/lofi/192.m3u8",
		Type: M3uRadio,
	},
	{
		Id:   "globo-sp",
		Name: "Globo SP",
		URL:  "https://medias.sgr.globo.com/hls/aRGloboSP/aRGloboSP.m3u8",
		Type: M3uRadio,
	},
	{
		Id:   "globo-rj",
		Name: "Globo RJ",
		URL:  "https://medias.sgr.globo.com/hls/aRGloboRJ/aRGloboRJ.m3u8",
		Type: M3uRadio,
	},
	{
		Id:   "swing",
		Name: "Swing songs (no ads)",
		URL:  "http://lainon.life:8000/swing.mp3",
		Type: M3uRadio,
	},
	{
		Id:   "cyber",
		Name: "Cyber songs (no ads)",
		URL:  "http://lainon.life:8000/cyberia.mp3",
		Type: M3uRadio,
	},
	{
		Id:   "cafe",
		Name: "Cafe songs (no ads)",
		URL:  "http://lainon.life:8000/cafe.mp3",
		Type: M3uRadio,
	},
}

func GetRadioList() []*RadioChannel {
	return radios
}

func GetRadioById(id string) *RadioChannel {
	for _, radio := range radios {
		if strings.EqualFold(id, radio.Id) {
			return radio
		}
	}
	return nil
}
