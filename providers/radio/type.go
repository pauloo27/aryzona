package radio

type RadioType struct {
	Name          string
	IsOppus       bool
	GetDirectURL  func(url string) string
	GetPlayingNow func(url, directURL string) (title, artist string)
}

type RadioChannel struct {
	ID, Name, URL string
	Type          *RadioType
}

func (c RadioChannel) CanPause() bool {
	return true
}

func (c RadioChannel) GetName() string {
	return c.Name + " (radio)"
}

func (c RadioChannel) GetDirectURL() (string, error) {
	return c.Type.GetDirectURL(c.URL), nil
}

func (c RadioChannel) GetFullTitle() (title, artist string) {
	directURL, err := c.GetDirectURL()
	if err != nil {
		return "", ""
	}
	return c.Type.GetPlayingNow(c.URL, directURL)
}

func (c RadioChannel) IsOppus() bool {
	return c.Type.IsOppus
}

func (c RadioChannel) IsLocal() bool {
	return false
}
