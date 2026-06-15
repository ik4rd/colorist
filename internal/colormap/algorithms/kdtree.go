package algorithms

import "github.com/ik4rd/colorist/internal/colormap"

func init() {
	colormap.Register(colormap.NewRecursive("kdtree", kdtreeSplit))
}

func kdtreeSplit(px *colormap.Pixels, r colormap.Rect, opts colormap.Options) []colormap.Rect {
	xAt, xCost, xOK := px.BestSplit(r, colormap.AxisVertical, opts.MinSize)
	yAt, yCost, yOK := px.BestSplit(r, colormap.AxisHorizontal, opts.MinSize)

	switch {
	case xOK && (!yOK || xCost <= yCost):
		return splitAt(r, colormap.AxisVertical, xAt)
	case yOK:
		return splitAt(r, colormap.AxisHorizontal, yAt)
	default:
		return nil
	}
}
