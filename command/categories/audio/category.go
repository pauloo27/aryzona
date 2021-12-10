package audio

import (
	"math/rand"
	"time"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/providers/youtube"
)

var Audio = command.CommandCategory{
	Name:  "Audio related stuff",
	Emoji: "🎵",
	Commands: []*command.Command{
		&LyricCommand, &RadioCommand, &PlayingCommand, &StopCommand, &PlayCommand,
		&SkipCommand, &PauseCommand,
	},
}

var (
	coolVids = []string{
		"Baby Shark",
		"Despacito",
		"Johny Jonhy",
		"Shape of you",
		"See you again",
		"Bath song",
		"Uptime funk open suse",
		"Dame To Cosita",
		"Sorry",
		"Lean On",
		"Gangnam Style",
	}
)

/* #nosec G404 */
func init() {
	command.RegisterCategory(Audio)
	// so, youtube ban apikeys with more than 90 days of inactivity,
	// since the bot is not used enought that i can asure that's never going to
	// happen, let's search something on every startup
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(coolVids))
	_, _, _ = youtube.GetBestResult(coolVids[index])
}
