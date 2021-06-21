package radio

import "strings"

type RadioType struct {
	Name         string
	IsOppus      bool
	GetDirectURL func(url string) string
}

type RadioChannel struct {
	Id, Name, URL string
	Type          *RadioType
}

var radios = []*RadioChannel{
	{
		Id:   "hunter-pop",
		Name: "Rádio Hunter POP",
		URL:  "https://hls.hunter.fm/pop/192.m3u8",
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
