package radio

import (
	"strings"
)

var radios = []*RadioChannel{
	{
		ID:   "hunter-pop",
		Name: "Rádio Hunter POP",
		URL:  "https://hls.hunter.fm/pop/192.m3u8",
		Type: HunterFM,
	},
	{
		ID:   "hunter-pisadinha",
		Name: "Rádio Hunter Pisadinha",
		URL:  "https://hls.hunter.fm/pisadinha/320.m3u8",
		Type: HunterFM,
	},
	{
		ID:   "hunter-tropical",
		Name: "Rádio Hunter Tropical",
		URL:  "https://hls.hunter.fm/tropical/192.m3u8",
		Type: HunterFM,
	},
	{
		ID:   "hunter-80s",
		Name: "Rádio Hunter Anos 80",
		URL:  "https://hls.hunter.fm/80s/192.m3u8",
		Type: HunterFM,
	},
	{
		ID:   "hunter-rock",
		Name: "Rádio Hunter Rock",
		URL:  "https://hls.hunter.fm/rock/192.m3u8",
		Type: HunterFM,
	},
	{
		ID:   "hunter-lofi",
		Name: "Rádio Hunter LoFi",
		URL:  "https://hls.hunter.fm/lofi/192.m3u8",
		Type: HunterFM,
	},
	{
		ID:   "cidade",
		Name: "Rádio Cidade",
		URL:  "https://18003.live.streamtheworld.com/RADIOCIDADEAAC.aac",
		Type: RadioCidade,
	},
	{
		ID:   "globo-sp",
		Name: "Globo SP",
		URL:  "https://medias.sgr.globo.com/hls/aRGloboSP/aRGloboSP.m3u8",
		Type: GenericRadio,
	},
	{
		ID:   "globo-rj",
		Name: "Globo RJ",
		URL:  "https://medias.sgr.globo.com/hls/aRGloboRJ/aRGloboRJ.m3u8",
		Type: GenericRadio,
	},
	{
		ID:   "swing",
		Name: "Swing songs (no ads)",
		URL:  "http://lainon.life:8000/swing.mp3",
		Type: LainchanRadio,
	},
	{
		ID:   "cyber",
		Name: "Cyber songs (no ads)",
		URL:  "http://lainon.life:8000/cyberia.mp3",
		Type: LainchanRadio,
	},
	{
		ID:   "cafe",
		Name: "Cafe songs (no ads)",
		URL:  "http://lainon.life:8000/cafe.mp3",
		Type: LainchanRadio,
	},
}

func GetRadioList() []*RadioChannel {
	return radios
}

func GetRadioByID(id string) *RadioChannel {
	for _, radio := range radios {
		if strings.EqualFold(id, radio.ID) {
			return radio
		}
	}
	return nil
}
