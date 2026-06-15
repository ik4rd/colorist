package colormap

import (
	"image/color"

	"github.com/ik4rd/colorist/internal/colormap/colornames"
)

type Region struct {
	X, Y, W, H int
	Mean       color.RGBA
	Hex        string
	RGB        string
	CMYK       string
	Name       string
}

func newRegion(r Rect, mean color.RGBA) Region {
	return Region{
		X:    r.X,
		Y:    r.Y,
		W:    r.W,
		H:    r.H,
		Mean: mean,
		Hex:  hex(mean),
		RGB:  rgbStr(mean),
		CMYK: cmykStr(mean),
		Name: colornames.Nearest(mean),
	}
}
