package f

import (
	"fmt"
	"strings"
	"time"

	"github.com/Pauloo27/aryzona/internal/i18n"
)

func Pluralize(i int, singular, plural string) string {
	if i == 1 {
		return singular
	}
	return plural
}

func PadRigth(source, padStr string, minLength int) string {
	if len(source) >= minLength {
		return source
	}
	return source + strings.Repeat(padStr, minLength-len(source))
}

func PadLeft(source, padStr string, minLength int) string {
	if len(source) >= minLength {
		return source
	}
	return strings.Repeat(padStr, minLength-len(source)) + source
}

func ShortDuration(duration time.Duration) string {
	durationSec := int(duration.Seconds())

	hours := durationSec / 3600
	minutes := (durationSec % 3600) / 60
	seconds := durationSec % 60

	if hours != 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func DurationAsText(duration time.Duration, t *i18n.Common) string {
	seconds := int(duration.Seconds())
	if seconds < 60 {
		return t.DurationLessThanAMinute.Str()
	}

	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	stringfy := func(i int, singular, plural string) string {
		if i == 0 {
			return ""
		}
		return fmt.Sprintf("%d %s", i, Pluralize(i, singular, plural))
	}

	return strings.TrimSpace(fmt.Sprintf(
		"%s %s %s",
		stringfy(days, t.DurationDay.Str(), t.DurationDays.Str()),
		stringfy(hours, t.DurationHour.Str(), t.DurationHours.Str()),
		stringfy(minutes, t.DurationMinute.Str(), t.DurationMinutes.Str()),
	))
}

func DurationAsDetailedDiffText(duration time.Duration, t *i18n.Common) string {
	totalSeconds := int(duration.Seconds())
	if totalSeconds == 0 {
		return t.DurationNow.Str()
	}

	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	stringfy := func(i int, singular, plural string) string {
		if i == 0 {
			return ""
		}
		return fmt.Sprintf("%d %s", i, Pluralize(i, singular, plural))
	}

	return strings.TrimSpace(fmt.Sprintf(
		"%s %s %s %s",
		stringfy(days, t.DurationDay.Str(), t.DurationDays.Str()),
		stringfy(hours, t.DurationHour.Str(), t.DurationHours.Str()),
		stringfy(minutes, t.DurationMinute.Str(), t.DurationMinutes.Str()),
		stringfy(seconds, t.DurationSecond.Str(), t.DurationSeconds.Str()),
	))
}

func Emojify(i int) string {
	if i < 10 {
		return numberEmojis[i]
	}
	sb := strings.Builder{}
	for _, c := range fmt.Sprintf("%d", i) {
		sb.WriteString(numberEmojis[int(c-'0')])
	}
	return sb.String()
}

var numberEmojis = map[int]string{
	0: "0️⃣",
	1: "1️⃣",
	2: "2️⃣",
	3: "3️⃣",
	4: "4️⃣",
	5: "5️⃣",
	6: "6️⃣",
	7: "7️⃣",
	8: "8️⃣",
	9: "9️⃣",
}
