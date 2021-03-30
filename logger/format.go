package logger

import (
	"fmt"
	"strconv"
)

const (
	ColorBold   = "\033[1m"
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorWhite  = "\033[39m"
)

func AsHexRGB(hex string) string {
	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)
	return AsRGB(r, g, b)
}

func AsRGB(r, g, b int64) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}
