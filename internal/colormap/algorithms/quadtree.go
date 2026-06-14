package algorithms

import "github.com/ik4rd/colorist/internal/colormap"

func init() {
	colormap.Register(colormap.NewRecursive("quadtree", quadtreeSplit))
}

const halvesPerAxis = 2

func quadtreeSplit(_ *colormap.Pixels, r colormap.Rect, opts colormap.Options) []colormap.Rect {
	if r.W < halvesPerAxis*opts.MinSize || r.H < halvesPerAxis*opts.MinSize {
		return nil
	}

	lw, lh := r.W/halvesPerAxis, r.H/halvesPerAxis
	rw, rh := r.W-lw, r.H-lh

	return []colormap.Rect{
		{X: r.X, Y: r.Y, W: lw, H: lh},
		{X: r.X + lw, Y: r.Y, W: rw, H: lh},
		{X: r.X, Y: r.Y + lh, W: lw, H: rh},
		{X: r.X + lw, Y: r.Y + lh, W: rw, H: rh},
	}
}
