package colormap

import (
	"bytes"
	"fmt"
	"image"

	"github.com/ik4rd/colorist/internal/imageio"
)

const DefaultAlgorithm = "quadtree"

const (
	LabelFormatNames = "names"
	LabelFormatHex   = "hex"
	LabelFormatRGB   = "rgb"
	LabelFormatCMYK  = "cmyk"
)

const (
	defaultThreshold     = 32
	defaultDetail        = 200
	defaultDetailFrac    = 0.05
	defaultMinSize       = 8
	defaultMaxDepth      = 64
	defaultGap           = 0
	defaultHalvesPerAxis = 2
	defaultLabelDensity  = 1.0
	defaultLabelFormat   = LabelFormatHex

	minBound         = 1
	minHalvesPerAxis = 2
)

type Options struct {
	Algorithm     string
	Threshold     float64
	Detail        float64
	DetailFrac    float64
	MinSize       int
	MaxDepth      int
	Gap           int
	HalvesPerAxis int
	LabelDensity  float64
	LabelFormat   string
}

func DefaultOptions() Options {
	return Options{
		Algorithm:     DefaultAlgorithm,
		Threshold:     defaultThreshold,
		Detail:        defaultDetail,
		DetailFrac:    defaultDetailFrac,
		MinSize:       defaultMinSize,
		MaxDepth:      defaultMaxDepth,
		Gap:           defaultGap,
		HalvesPerAxis: defaultHalvesPerAxis,
		LabelDensity:  defaultLabelDensity,
		LabelFormat:   defaultLabelFormat,
	}
}

func (o Options) normalized() Options {
	if o.MinSize < minBound {
		o.MinSize = minBound
	}
	if o.MaxDepth < minBound {
		o.MaxDepth = minBound
	}
	if o.HalvesPerAxis < minHalvesPerAxis {
		o.HalvesPerAxis = minHalvesPerAxis
	}
	if o.Algorithm == "" {
		o.Algorithm = DefaultAlgorithm
	}

	return o
}

func NewPixels(img image.Image) (*Pixels, error) {
	if img == nil {
		return nil, fmt.Errorf("colormap: nil image")
	}

	b := img.Bounds()
	if b.Dx() == 0 || b.Dy() == 0 {
		return nil, fmt.Errorf("colormap: empty image")
	}

	return newPixels(img), nil
}

func BuildFrom(px *Pixels, opts Options) ([]Region, error) {
	if px == nil {
		return nil, fmt.Errorf("colormap: nil pixels")
	}

	opts = opts.normalized()

	algo, err := getAlgorithm(opts.Algorithm)
	if err != nil {
		return nil, err
	}

	return algo.Partition(px, opts), nil
}

func Build(img image.Image, opts Options) ([]Region, error) {
	px, err := NewPixels(img)
	if err != nil {
		return nil, err
	}

	return BuildFrom(px, opts)
}

func ProcessPixels(px *Pixels, opts Options) (image.Image, error) {
	return ProcessPixelsView(px, View{}, opts)
}

func ProcessPixelsView(px *Pixels, view View, opts Options) (image.Image, error) {
	regions, err := BuildFrom(px, opts)
	if err != nil {
		return nil, err
	}

	return RenderView(regions, px.w, px.h, view, opts), nil
}

func Process(img image.Image, opts Options) (image.Image, error) {
	px, err := NewPixels(img)
	if err != nil {
		return nil, err
	}

	return ProcessPixels(px, opts)
}

func ProcessBytes(img image.Image, format string, opts Options) ([]byte, error) {
	out, err := Process(img, opts)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := imageio.Encode(&buf, format, out); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
