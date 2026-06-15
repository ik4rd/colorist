package algorithms

import "github.com/ik4rd/colorist/internal/colormap"

func init() {
	colormap.Register(colormap.NewRecursive("bsp", bspSplit))
}

func bspSplit(_ *colormap.Pixels, r colormap.Rect, opts colormap.Options) []colormap.Rect {
	n := max(opts.HalvesPerAxis, 2)

	if r.W >= r.H {
		if r.W < n*opts.MinSize {
			return nil
		}
		xs := splitAxis(r.X, r.W, n)
		rects := make([]colormap.Rect, n)
		for i := range n {
			rects[i] = colormap.Rect{X: xs[i], Y: r.Y, W: xs[i+1] - xs[i], H: r.H}
		}
		return rects
	}

	if r.H < n*opts.MinSize {
		return nil
	}
	ys := splitAxis(r.Y, r.H, n)
	rects := make([]colormap.Rect, n)
	for i := range n {
		rects[i] = colormap.Rect{X: r.X, Y: ys[i], W: r.W, H: ys[i+1] - ys[i]}
	}
	return rects
}
