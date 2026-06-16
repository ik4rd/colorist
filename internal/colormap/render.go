package colormap

import (
	"image"
	"image/draw"
	"math"
	"sort"
)

const (
	renderTarget = 2000
	cropMaxOut   = 2560
)

func (v View) empty() bool { return v.W <= 0 || v.H <= 0 }

func Render(regions []Region, w, h int, opts Options) image.Image {
	return RenderView(regions, w, h, View{}, opts)
}

func RenderView(regions []Region, w, h int, view View, opts Options) image.Image {
	S := renderScale(w, h)
	base := float64(max(w, h))

	var ox, oy float64
	var outW, outH int
	var g float64
	var lblScale float64

	if view.empty() {
		ox, oy = 0, 0
		g = float64(S)
		outW, outH = w*S, h*S
		lblScale = float64(S)
	} else {
		ox, oy = float64(view.X), float64(view.Y)
		out := math.Min(base*float64(S), cropMaxOut)
		g = out / float64(max(view.W, view.H))
		outW = int(math.Round(float64(view.W) * g))
		outH = int(math.Round(float64(view.H) * g))
		lblScale = out / base
	}

	dst := image.NewRGBA(image.Rect(0, 0, outW, outH))
	bounds := dst.Bounds()
	lbl := newLabeler(lblScale)

	type lab struct {
		rect               image.Rectangle
		primary, secondary string
		ref                string
		skip               bool
	}

	labs := make([]lab, len(regions))
	eligible := make([]int, 0, len(regions))

	gap := float64(opts.Gap)

	for i, rg := range regions {
		fx0 := (float64(rg.X) - ox) * g
		fy0 := (float64(rg.Y) - oy) * g
		fx1 := (float64(rg.X+rg.W) - ox) * g
		fy1 := (float64(rg.Y+rg.H) - oy) * g

		if gap > 0 {
			fx1 -= gap * g
			fy1 -= gap * g
		}

		rect := image.Rect(
			int(math.Round(fx0)), int(math.Round(fy0)),
			int(math.Round(fx1)), int(math.Round(fy1)),
		).Intersect(bounds)

		if rect.Empty() {
			labs[i].skip = true
			continue
		}

		primary, secondary, ref := labelLines(rg, opts.LabelFormat)
		labs[i] = lab{rect: rect, primary: primary, secondary: secondary, ref: ref}

		if lbl.fits(rect, ref) {
			eligible = append(eligible, i)
		}
	}

	keep := selectLabeled(regions, eligible, opts.LabelDensity)

	for i, rg := range regions {
		l := labs[i]
		if l.skip {
			continue
		}

		draw.Draw(dst, l.rect, &image.Uniform{C: rg.Mean}, image.Point{}, draw.Src)

		if keep[i] {
			lbl.draw(dst, l.rect, l.primary, l.secondary, l.ref, rg.Mean)
		}
	}

	return dst
}

func renderScale(w, h int) int {
	longest := max(h, w)

	if longest <= 0 {
		return 1
	}

	scale := min(max((renderTarget+longest-1)/longest, 1), 4)

	return scale
}

func selectLabeled(regions []Region, eligible []int, density float64) []bool {
	keep := make([]bool, len(regions))
	if density <= 0 || len(eligible) == 0 {
		return keep
	}

	if density >= 1 {
		for _, i := range eligible {
			keep[i] = true
		}
		return keep
	}

	n := int(math.Ceil(float64(len(eligible)) * density))
	if n <= 0 {
		return keep
	}

	sort.SliceStable(eligible, func(a, b int) bool {
		ra, rb := regions[eligible[a]], regions[eligible[b]]
		return ra.W*ra.H > rb.W*rb.H
	})

	for _, i := range eligible[:n] {
		keep[i] = true
	}

	return keep
}

func labelLines(rg Region, format string) (primary, secondary, ref string) {
	widest := rg.Hex
	if len(rg.RGB) > len(widest) {
		widest = rg.RGB
	}
	if len(rg.CMYK) > len(widest) {
		widest = rg.CMYK
	}

	switch format {
	case LabelFormatNames:
		return rg.Name, "", rg.Name
	case LabelFormatRGB:
		return rg.RGB, rg.Name, widest
	case LabelFormatCMYK:
		return rg.CMYK, rg.Name, widest
	default:
		return rg.Hex, rg.Name, widest
	}
}
