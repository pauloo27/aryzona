package model

type EmbedField struct {
	Name, Value string
	Inline      bool
}

type Embed struct {
	Fields       []*EmbedField
	Title        string
	Description  string
	ImageURL     string
	ThumbnailURL string
	Footer       string
	Color        int
	URL          string
}

func NewEmbed() *Embed {
	return &Embed{}
}

func (e *Embed) WithColor(color int) *Embed {
	e.Color = color
	return e
}

func (e *Embed) WithTitle(title string) *Embed {
	e.Title = title
	return e
}

func (e *Embed) WithDescription(description string) *Embed {
	e.Description = description
	return e
}

func (e *Embed) WithFieldInline(name, value string) *Embed {
	e.Fields = append(e.Fields, &EmbedField{
		Name:   name,
		Value:  value,
		Inline: true,
	})
	return e
}

func (e *Embed) WithField(name, value string) *Embed {
	e.Fields = append(e.Fields, &EmbedField{
		Name:   name,
		Value:  value,
		Inline: false,
	})
	return e
}

func (e *Embed) WithImage(url string) *Embed {
	e.ImageURL = url
	return e
}

func (e *Embed) WithThumbnail(url string) *Embed {
	e.ThumbnailURL = url
	return e
}

func (e *Embed) WithURL(url string) *Embed {
	e.URL = url
	return e
}

func (e *Embed) WithFooter(text string) *Embed {
	e.Footer = text
	return e
}
