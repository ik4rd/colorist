package colormap

type Rect struct {
	X, Y, W, H int
}

type SplitFunc func(px *Pixels, r Rect, opts Options) []Rect

func Subdivide(px *Pixels, opts Options, split SplitFunc) []Region {
	var regions []Region

	var walk func(r Rect, depth int)
	walk = func(r Rect, depth int) {
		mean, refine := px.heterogeneous(r, opts)

		if refine && depth < opts.MaxDepth {
			if children := split(px, r, opts); len(children) > 0 {
				for _, c := range children {
					walk(c, depth+1)
				}
				return
			}
		}

		regions = append(regions, newRegion(r, mean))
	}

	walk(Rect{X: 0, Y: 0, W: px.w, H: px.h}, 0)

	return regions
}

type recursiveAlgorithm struct {
	name  string
	split SplitFunc
}

func (a recursiveAlgorithm) Name() string { return a.name }

func (a recursiveAlgorithm) Partition(px *Pixels, opts Options) []Region {
	return Subdivide(px, opts, a.split)
}

func NewRecursive(name string, split SplitFunc) Algorithm {
	return recursiveAlgorithm{name: name, split: split}
}
