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
		return fmt.Sprintf("%d %s", i, Pluralize(i, singular, plural))
	}

	return strings.TrimSpace(fmt.Sprintf(
		"%s %s %s",
		stringfy(days, "day", "days"),
		stringfy(hours, "hour", "hours"),
		stringfy(minutes, "minute", "minutes"),
	))
}
