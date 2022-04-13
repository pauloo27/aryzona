package radio

import (
	"io"
	"net/http"
	"regexp"
)

type CidadeRadio struct {
	BaseRadio
	ID, Name, URL string
}

var _ RadioChannel = &CidadeRadio{}

func newCidadeRadio(id, name, url string) CidadeRadio {
	return CidadeRadio{
		ID:        id,
		Name:      name,
		URL:       url,
		BaseRadio: BaseRadio{},
	}
}

func (r CidadeRadio) GetID() string {
	return r.ID
}

func (r CidadeRadio) GetName() string {
	return r.Name
}

func (r CidadeRadio) GetShareURL() string {
	return "https://radiocidade.fm/"
}

func (r CidadeRadio) GetThumbnailURL() (string, error) {
	return "", nil
}

func (r CidadeRadio) IsOpus() bool {
	return false
}

func (r CidadeRadio) GetDirectURL() (string, error) {
	return r.URL, nil
}

func (r CidadeRadio) GetFullTitle() (title, artist string) {
	res, err := http.Get("https://np.tritondigital.com/public/nowplaying?mountName=RADIOCIDADEAAC&numberToFetch=1&eventType=track")
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	// nobody deserves to deal with XML... lets just pretend that we got a string
	// as response and a regex is the way to parse it...
	bodyStr := string(body)
	parseRegex := regexp.MustCompile(`CDATA\[([^\]]+)\]`)
	matches := parseRegex.FindAllStringSubmatch(bodyStr, -1)
	if len(matches) < 4 {
		return
	}
	title = matches[2][1]
	artist = matches[3][1]
	return
}
