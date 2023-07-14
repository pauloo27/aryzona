package i18n

import (
	"time"

	"github.com/goodsign/monday"
)

type Locale struct {
	langName LanguageName

	SimpleDateTimeFormat string
}

func (l *Locale) FormatSimpleDateTime(time time.Time) string {
	a := monday.Format(
		time, l.SimpleDateTimeFormat, monday.Locale(l.langName),
	)

	return a
}
