package colormap

import "image/color"

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

type Region struct {
	X, Y, W, H int
	Mean       color.RGBA
	Hex        string
	RGB        string
	CMYK       string
	Name       string
}

type Rect struct {
	X, Y, W, H int
}

type View struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Pixels struct {
	w, h int
	r    []float64
	g    []float64
	b    []float64
}
