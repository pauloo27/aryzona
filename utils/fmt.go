package utils

import (
	"fmt"
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

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	var out string

	if minutes > 0 {
		out = Fmt("%d %s", minutes, Pluralize(minutes, "minute", "minutes"))
	}

	if hours > 0 {
		out = Fmt("%d %s and %s", hours, Pluralize(hours, "hour", "hours"), out)
	}

	return out
}
