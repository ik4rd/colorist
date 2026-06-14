package colormap

import (
	"image"
	"image/draw"
)

func Render(regions []Region, w, h int, opts Options) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for _, rg := range regions {
		x0, y0 := rg.X, rg.Y
		x1, y1 := rg.X+rg.W, rg.Y+rg.H

		if opts.Gap > 0 {
			x1 -= opts.Gap
			y1 -= opts.Gap
			if x1 <= x0 || y1 <= y0 {
				continue
			}
		}

		draw.Draw(dst, image.Rect(x0, y0, x1, y1), &image.Uniform{C: rg.Mean}, image.Point{}, draw.Src)
	}

	return dst
}
