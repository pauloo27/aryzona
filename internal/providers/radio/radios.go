package radio

var (
	radioMap  = make(map[string]RadioChannel)
	radioList []RadioChannel
)

func init() {
	registerRadios(
		newYouTubeRadio("lofi", "Lofi: beats to relax/study", "https://youtube.com/watch?v=5qap5aO4i9A"),
		newYouTubeRadio("techno", "Techno: Rave Radio 24/7", "https://www.youtube.com/watch?v=6Irus3d5f0E"),

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
			"lofi2", "Rádio Hunter Lofi", "https://hls.hunter.fm/lofi/192.m3u8",
		),

		newGloboRadio(
			"globo-rj", "Rádio Globo RJ", "https://medias.sgr.globo.com/hls/aRGloboRJ/aRGloboRJ.m3u8",
		),
		newGloboRadio(
			"globo-sp", "Rádio Globo SP", "https://medias.sgr.globo.com/hls/aRGloboSP/aRGloboSP.m3u8",
		),

		newLainchanRadio("cyber", "Cyber songs (no ads)", "http://lainon.life:8000/mpd-cyberia.mp3"),
		newLainchanRadio("cafe", "Cafe songs (no ads)", "http://lainon.life:8000/mpd-cafe.mp3"),
		newLainchanRadio("swing", "Swing songs (no ads)", "http://lainon.life:8000/mpd-swing.mp3"),
		// i dont know what to name that one...
		newLainchanRadio("all", "Everything songs (no ads)", "http://lainon.life:8000/mpd-everything.mp3"),
	)
}

func registerRadios(radios ...RadioChannel) {
	for _, radio := range radios {
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
