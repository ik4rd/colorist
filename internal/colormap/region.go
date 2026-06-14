package colormap

import "image/color"

type Region struct {
	X, Y, W, H int
	Mean       color.RGBA
	Hex        string
}

func newRegion(r Rect, mean color.RGBA) Region {
	return Region{X: r.X, Y: r.Y, W: r.W, H: r.H, Mean: mean, Hex: hex(mean)}
}
