package algorithms

import "github.com/ik4rd/colorist/internal/colormap"

func init() {
	colormap.Register(colormap.NewRecursive("adaptive_bsp", adaptiveSplit))
}

func adaptiveSplit(px *colormap.Pixels, r colormap.Rect, opts colormap.Options) []colormap.Rect {
	axis := colormap.AxisVertical
	if r.H > r.W {
		axis = colormap.AxisHorizontal
	}

	at, _, ok := px.BestSplit(r, axis, opts.MinSize)
	if !ok {
		return nil
	}

	return splitAt(r, axis, at)
}

func splitAt(r colormap.Rect, axis colormap.Axis, at int) []colormap.Rect {
	if axis == colormap.AxisVertical {
		return []colormap.Rect{
			{X: r.X, Y: r.Y, W: at - r.X, H: r.H},
			{X: at, Y: r.Y, W: r.X + r.W - at, H: r.H},
		}
	}

	return []colormap.Rect{
		{X: r.X, Y: r.Y, W: r.W, H: at - r.Y},
		{X: r.X, Y: at, W: r.W, H: r.Y + r.H - at},
	}
}
