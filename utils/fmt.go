package utils

import (
	"fmt"
	"strings"
	"time"
)

func Pluralize(i int, singular, plural string) string {
	if i == 1 {
		return singular
	}
	return plural
}

func Fmt(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

func FormatDuration(duration time.Duration) string {
	seconds := int(duration.Seconds())
	if seconds < 60 {
		return "less than a minute"
	}

	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	stringfy := func(i int, singular, plural string) string {
		if i == 0 {
			return ""
		}
		return Fmt("%d %s", i, Pluralize(i, singular, plural))
	}

	return strings.TrimSpace(Fmt(
		"%s %s %s",
		stringfy(days, "day", "days"),
		stringfy(hours, "hour", "hours"),
		stringfy(minutes, "minute", "minutes"),
	))
}
