package radio

var M3uRadio = &RadioType{
	Name:    "m3ue playlist",
	IsOppus: false,
	GetDirectURL: func(url string) string {
		return url
	},
}
