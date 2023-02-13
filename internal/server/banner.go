package server

import (
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"github.com/Pauloo27/aryzona/internal/providers/livescore"
	"github.com/go-chi/chi/v5"
)

const (
	teamImgHeight = 64
	imgPadding    = 10
	imgHeight     = teamImgHeight + imgPadding*2
	imgWidth      = imgHeight * 3
)

var (
	vsImg image.Image
)

func renderBanner(w http.ResponseWriter, r *http.Request) {
	if vsImg == nil {
		err := loadVsImg()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	m := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	t1Img, err := loadTeamImg(chi.URLParam(r, "t1"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t2Img, err := loadTeamImg(chi.URLParam(r, "t2"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	draw.Draw(
		m, m.Bounds(),
		t1Img,
		image.Point{-imgPadding, -imgPadding},
		draw.Over,
	)
	draw.Draw(
		m, m.Bounds(),
		t2Img,
		image.Point{-(imgWidth - imgPadding - teamImgHeight), -imgPadding},
		draw.Over,
	)

	draw.Draw(
		m, m.Bounds(),
		vsImg,
		image.Point{-((imgWidth - teamImgHeight) / 2), -imgPadding},
		draw.Over,
	)

	err = png.Encode(w, m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/* #nosec GG107 */
func loadTeamImg(teamID string) (image.Image, error) {
	url := livescore.GetTeamImgURL(teamID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	/* #nosec G307 */
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	return img, err
}

func loadVsImg() error {
	f, err := os.OpenFile("assets/vs.png", os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(f)
	vsImg = img
	return err
}
