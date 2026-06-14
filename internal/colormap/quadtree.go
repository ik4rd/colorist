package colormap

func init() {
	Register(quadtree{})
}

const halvesPerAxis = 2

type quadtree struct{}

func (quadtree) Name() string { return "quadtree" }

func (quadtree) Partition(px *pixels, opts Options) []Region {
	var regions []Region

	var split func(x, y, w, h, depth int)
	split = func(x, y, w, h, depth int) {
		mean, refine := px.heterogeneous(x, y, w, h, opts)

		divisible := w >= halvesPerAxis*opts.MinSize && h >= halvesPerAxis*opts.MinSize
		if !refine || depth >= opts.MaxDepth || !divisible {
			regions = append(regions, newRegion(x, y, w, h, mean))
			return
		}

		lw, lh := w/halvesPerAxis, h/halvesPerAxis
		rw, rh := w-lw, h-lh

		split(x, y, lw, lh, depth+1)
		split(x+lw, y, rw, lh, depth+1)
		split(x, y+lh, lw, rh, depth+1)
		split(x+lw, y+lh, rw, rh, depth+1)
	}

	split(0, 0, px.w, px.h, 0)

	return regions
}
