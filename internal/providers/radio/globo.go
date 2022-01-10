package radio

type GloboRadio struct {
	BaseRadio
	ID, Name, URL string
}

var _ RadioChannel = &GloboRadio{}

func newGloboRadio(id, name, url string) GloboRadio {
	return GloboRadio{
		ID:        id,
		Name:      name,
		URL:       url,
		BaseRadio: BaseRadio{},
	}
}

func (r GloboRadio) GetID() string {
	return r.ID
}

func (r GloboRadio) GetName() string {
	return r.Name
}

func (r GloboRadio) GetThumbnailURL() (string, error) {
	return "", nil
}

func (r GloboRadio) IsOpus() bool {
	return false
}

func (r GloboRadio) GetDirectURL() (string, error) {
	return r.URL, nil
}

func (r GloboRadio) GetFullTitle() (title, artist string) {
	return r.GetName(), ""
}
