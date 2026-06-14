package colormap

import "image/color"

type Region struct {
	X, Y, W, H int
	Mean       color.RGBA
	Hex        string
}

func newRegion(x, y, w, h int, mean color.RGBA) Region {
	return Region{X: x, Y: y, W: w, H: h, Mean: mean, Hex: hex(mean)}
}
