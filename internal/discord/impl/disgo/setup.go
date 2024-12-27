package disgo

import "github.com/pauloo27/aryzona/internal/discord"

func init() {
	discord.UseImplementation(&DisgoBot{})
}
