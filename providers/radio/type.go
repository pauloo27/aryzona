package radio

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
