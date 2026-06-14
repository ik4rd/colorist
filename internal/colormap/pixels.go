package colormap

import (
	"image"
	"image/color"
	"math"
)

const (
	channelShift = 8
	channelCount = 3
)

type pixels struct {
	w, h int
	r    []float64
	g    []float64
	b    []float64
}

func newPixels(img image.Image) *pixels {
	bnd := img.Bounds()
	w, h := bnd.Dx(), bnd.Dy()

	p := &pixels{w: w, h: h, r: make([]float64, w*h), g: make([]float64, w*h), b: make([]float64, w*h)}

	for y := range h {
		for x := range w {
			cr, cg, cb, _ := img.At(bnd.Min.X+x, bnd.Min.Y+y).RGBA()
			i := y*w + x
			p.r[i] = float64(cr >> channelShift)
			p.g[i] = float64(cg >> channelShift)
			p.b[i] = float64(cb >> channelShift)
		}
	}

	return p
}

func (p *pixels) heterogeneous(x, y, w, h int, opts Options) (color.RGBA, bool) {
	mean, rms := p.stats(x, y, w, h)
	if rms > opts.Threshold {
		return mean, true
	}

	if opts.Detail > 0 && opts.DetailFrac > 0 &&
		p.outlierFraction(x, y, w, h, mean, opts.Detail) > opts.DetailFrac {
		return mean, true
	}

	return mean, false
}

func (p *pixels) stats(x, y, w, h int) (mean color.RGBA, rms float64) {
	var sr, sg, sb, sr2, sg2, sb2 float64

	n := float64(w * h)

	for j := y; j < y+h; j++ {
		row := j * p.w
		for i := x; i < x+w; i++ {
			k := row + i
			vr, vg, vb := p.r[k], p.g[k], p.b[k]
			sr += vr
			sg += vg
			sb += vb
			sr2 += vr * vr
			sg2 += vg * vg
			sb2 += vb * vb
		}
	}

	mr, mg, mb := sr/n, sg/n, sb/n

	varR := math.Max(0, sr2/n-mr*mr)
	varG := math.Max(0, sg2/n-mg*mg)
	varB := math.Max(0, sb2/n-mb*mb)

	rms = math.Sqrt((varR + varG + varB) / channelCount)
	mean = color.RGBA{R: clamp8(mr), G: clamp8(mg), B: clamp8(mb), A: maxChannel}

	return mean, rms
}

func (p *pixels) outlierFraction(x, y, w, h int, mean color.RGBA, dist float64) float64 {
	mr, mg, mb := float64(mean.R), float64(mean.G), float64(mean.B)
	thr := dist * dist

	var outliers int

	for j := y; j < y+h; j++ {
		row := j * p.w
		for i := x; i < x+w; i++ {
			k := row + i
			dr, dg, db := p.r[k]-mr, p.g[k]-mg, p.b[k]-mb
			if dr*dr+dg*dg+db*db > thr {
				outliers++
			}
		}
	}

	return float64(outliers) / float64(w*h)
}
