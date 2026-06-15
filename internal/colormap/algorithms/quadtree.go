package algorithms

import "github.com/ik4rd/colorist/internal/colormap"

func init() {
	colormap.Register(colormap.NewRecursive("quadtree", quadtreeSplit))
}

func quadtreeSplit(_ *colormap.Pixels, r colormap.Rect, opts colormap.Options) []colormap.Rect {
	n := max(opts.HalvesPerAxis, 2)

	if r.W < n*opts.MinSize || r.H < n*opts.MinSize {
		return nil
	}

	xs := splitAxis(r.X, r.W, n)
	ys := splitAxis(r.Y, r.H, n)

	rects := make([]colormap.Rect, 0, n*n)

	for j := range n {
		for i := range n {
			rects = append(rects, colormap.Rect{
				X: xs[i],
				Y: ys[j],
				W: xs[i+1] - xs[i],
				H: ys[j+1] - ys[j],
			})
		}
	}

	return rects
}

func splitAxis(start, length, n int) []int {
	bounds := make([]int, n+1)
	for i := 0; i <= n; i++ {
		bounds[i] = start + i*length/n
	}

	return bounds
}
