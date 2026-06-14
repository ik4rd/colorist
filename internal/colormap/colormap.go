package colormap

import (
	"fmt"
	"image"
)

const DefaultAlgorithm = "quadtree"

const (
	defaultThreshold  = 32
	defaultDetail     = 200
	defaultDetailFrac = 0.05
	defaultMinSize    = 8
	defaultMaxDepth   = 64
	defaultGap        = 0

	minBound = 1
)

type Options struct {
	Algorithm  string
	Threshold  float64
	Detail     float64
	DetailFrac float64
	MinSize    int
	MaxDepth   int
	Gap        int
}

func DefaultOptions() Options {
	return Options{
		Algorithm:  DefaultAlgorithm,
		Threshold:  defaultThreshold,
		Detail:     defaultDetail,
		DetailFrac: defaultDetailFrac,
		MinSize:    defaultMinSize,
		MaxDepth:   defaultMaxDepth,
		Gap:        defaultGap,
	}
}

func Build(img image.Image, opts Options) ([]Region, error) {
	if img == nil {
		return nil, fmt.Errorf("colormap: nil image")
	}

	if opts.MinSize < minBound {
		opts.MinSize = minBound
	}
	if opts.MaxDepth < minBound {
		opts.MaxDepth = minBound
	}
	if opts.Algorithm == "" {
		opts.Algorithm = DefaultAlgorithm
	}

	b := img.Bounds()
	if b.Dx() == 0 || b.Dy() == 0 {
		return nil, fmt.Errorf("colormap: empty image")
	}

	algo, err := getAlgorithm(opts.Algorithm)
	if err != nil {
		return nil, err
	}

	px := newPixels(img)

	return algo.Partition(px, opts), nil
}
