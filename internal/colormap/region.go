package colormap

import (
	"image/color"

	"github.com/ik4rd/colorist/internal/colormap/colornames"
)

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
