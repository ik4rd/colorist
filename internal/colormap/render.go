package colormap

import (
	"image"
	"image/draw"
	"math"
	"sort"
)

const renderTarget = 2000

func Render(regions []Region, w, h int, opts Options) image.Image {
	scale := renderScale(w, h)
	dst := image.NewRGBA(image.Rect(0, 0, w*scale, h*scale))

	lbl := newLabeler(scale)

	type lab struct {
		rect               image.Rectangle
		primary, secondary string
		ref                string
		skip               bool
	}

	labs := make([]lab, len(regions))
	eligible := make([]int, 0, len(regions))

	for i, rg := range regions {
		x0, y0 := rg.X*scale, rg.Y*scale
		x1, y1 := (rg.X+rg.W)*scale, (rg.Y+rg.H)*scale

		if opts.Gap > 0 {
			x1 -= opts.Gap * scale
			y1 -= opts.Gap * scale
			if x1 <= x0 || y1 <= y0 {
				labs[i].skip = true
				continue
			}
		}

		rect := image.Rect(x0, y0, x1, y1)
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
