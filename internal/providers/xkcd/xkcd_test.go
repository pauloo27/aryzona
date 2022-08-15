package xkcd_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/providers/xkcd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByNum(t *testing.T) {
	// yeah so if there's a change in the response, the test needs to be changed.

	t.Run("get the first one", func(t *testing.T) {
		comic, err := xkcd.GetByNum(1)
		require.Nil(t, err)
		assert.Equal(t, 1, comic.Num)
		assert.Equal(t, "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg", comic.Img)
		assert.Equal(t, "Barrel - Part 1", comic.Title)
		assert.Equal(t, "Don't we all.", comic.Alt)
		assert.Equal(t, "2006", comic.Year)
		assert.Equal(t, "1", comic.Month)
		assert.Equal(t, "1", comic.Day)
		assert.Equal(t, "[[A boy sits in a barrel which is floating in an ocean.]]\nBoy: I wonder where I'll float next?\n[[The barrel drifts into the distance. Nothing else can be seen.]]\n{{Alt: Don't we all.}}", comic.Transcript)
		assert.Equal(t, "Barrel - Part 1", comic.SafeTitle)
		assert.Equal(t, "", comic.News)
	})

	t.Run("get april fools one", func(t *testing.T) {
		comic, err := xkcd.GetByNum(2601)
		require.Nil(t, err)
		assert.Equal(t, 2601, comic.Num)
		assert.Equal(t, "https://imgs.xkcd.com/comics/instructions.png", comic.Img)
		assert.Equal(t, "Instructions", comic.Title)
		assert.Equal(t, "Happy little turtles", comic.Alt)
		assert.Equal(t, "2022", comic.Year)
		assert.Equal(t, "4", comic.Month)
		assert.Equal(t, "1", comic.Day)
		assert.Equal(t, "", comic.Transcript)
		assert.Equal(t, "Instructions", comic.SafeTitle)
		assert.Equal(t, `Today's comic was created with <a href="https://instagram.com/fading_interest">Patrick</a>, <a href="https://twitter.com/Aiiane">Amber</a>, <a href="https://twitter.com/chromakode">@chromakode</a>, <a href="https://twitter.com/dyfrgi">Michael</a>, <a href="https://twitter.com/wirehead2501">Kat</a>, <a href="https://twitter.com/xDirtyPunkx">Conor</a>, <a href="https://twitter.com/zigdon">@zigdon</a>, and <a href="https://twitter.com/bstaffin">Benjamin Staffin</a>.`, comic.News)
	})
}

func TestRandom(t *testing.T) {
	t.Run("get a random one", func(t *testing.T) {
		comic, err := xkcd.GetRandom()
		require.Nil(t, err)
		// maybe there's a comic without a title or a img, which might be not a
		// good thing (if you are me from the future coming here to fix the test,
		// remember to fix the GetLatest test too cuz i used the same stupid asserts).
		assert.NotEqual(t, 0, comic.Num)
		assert.NotEqual(t, "", comic.Img)
		assert.NotEqual(t, "", comic.Title)
	})
}

func TestLatest(t *testing.T) {
	t.Run("get the latest one", func(t *testing.T) {
		comic, err := xkcd.GetLatest()
		require.Nil(t, err)
		assert.NotEqual(t, 0, comic.Num)
		assert.NotEqual(t, "", comic.Img)
		assert.NotEqual(t, "", comic.Title)
	})
}
