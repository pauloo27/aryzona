package radio

var (
	radioMap  = make(map[string]RadioChannel)
	radioList []RadioChannel
)

func init() {
	registerRadios(
		newYouTubeRadio("lofi", "Lofi: beats to relax/study to", "https://www.youtube.com/watch?v=jfKfPfyJRdk"),
		newYouTubeRadio("lofi-sleep", "Lofi: beats to sleep/chill to", "https://www.youtube.com/watch?v=rUxyKA_-grg"),

		newCidadeRadio(
			"cidade", "Rádio Cidade", "https://18003.live.streamtheworld.com/RADIOCIDADEAAC.aac",
		),

		newHunterRadio(
			"pisadinha", "Rádio Hunter Pisadinha", "https://hls.hunter.fm/pisadinha/320.m3u8",
		),
		newHunterRadio(
			"pop", "Rádio Hunter Pop", "https://hls.hunter.fm/pop/192.m3u8",
		),
		newHunterRadio(
			"pop2k", "Rádio Hunter Pop 2k", "https://hls.hunter.fm/pop2k/192.m3u8",
		),
		newHunterRadio(
			"rock", "Rádio Hunter Rock", "https://hls.hunter.fm/rock/192.m3u8",
		),
		newHunterRadio(
			"sertanejo", "Rádio Hunter Sertanejo", "https://hls.hunter.fm/sertanejo/192.m3u8",
		),
		newHunterRadio(
			"smash", "Rádio Hunter Smash", "https://hls.hunter.fm/smash/192.m3u8",
		),
		newHunterRadio(
			"80s", "Rádio Hunter 80s", "https://hls.hunter.fm/80s/192.m3u8",
		),
		newHunterRadio(
			"tropical", "Rádio Hunter Tropical", "https://hls.hunter.fm/tropical/192.m3u8",
		),
		newHunterRadio(
			"lofi-hunter", "Rádio Hunter Lofi", "https://hls.hunter.fm/lofi/192.m3u8",
		),
	)
}

func registerRadios(radios ...RadioChannel) {
	for _, radio := range radios {
		if radio == nil {
			continue
		}
		radioList = append(radioList, radio)
		radioMap[radio.GetID()] = radio
	}
}

func GetRadioList() []RadioChannel {
	return radioList
}

func GetRadioByID(id string) RadioChannel {
	return radioMap[id]
}
