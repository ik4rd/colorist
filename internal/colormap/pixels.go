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

func newPixels(img image.Image) *Pixels {
	bnd := img.Bounds()
	w, h := bnd.Dx(), bnd.Dy()

	p := &Pixels{w: w, h: h, r: make([]float64, w*h), g: make([]float64, w*h), b: make([]float64, w*h)}

	switch src := img.(type) {
	case *image.YCbCr:
		for y := range h {
			for x := range w {
				yi := src.YOffset(bnd.Min.X+x, bnd.Min.Y+y)
				ci := src.COffset(bnd.Min.X+x, bnd.Min.Y+y)
				cr, cg, cb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				p.set(y*w+x, cr, cg, cb)
			}
		}
	case *image.RGBA:
		for y := range h {
			o := src.PixOffset(bnd.Min.X, bnd.Min.Y+y)
			for x := range w {
				p.set(y*w+x, src.Pix[o], src.Pix[o+1], src.Pix[o+2])
				o += 4
			}
		}
	case *image.NRGBA:
		for y := range h {
			o := src.PixOffset(bnd.Min.X, bnd.Min.Y+y)
			for x := range w {
				a := uint32(src.Pix[o+3])
				p.set(y*w+x,
					uint8(uint32(src.Pix[o])*a/maxChannel),
					uint8(uint32(src.Pix[o+1])*a/maxChannel),
					uint8(uint32(src.Pix[o+2])*a/maxChannel))
				o += 4
			}
		}
	case *image.Gray:
		for y := range h {
			o := src.PixOffset(bnd.Min.X, bnd.Min.Y+y)
			for x := range w {
				v := src.Pix[o]
				p.set(y*w+x, v, v, v)
				o++
			}
		}
	default:
		for y := range h {
			for x := range w {
				cr, cg, cb, _ := img.At(bnd.Min.X+x, bnd.Min.Y+y).RGBA()
				p.set(y*w+x, uint8(cr>>channelShift), uint8(cg>>channelShift), uint8(cb>>channelShift))
			}
		}
	}

	return p
}

func (p *Pixels) set(i int, r, g, b uint8) {
	p.r[i] = float64(r)
	p.g[i] = float64(g)
	p.b[i] = float64(b)
}

func (p *Pixels) Stats(r Rect) (mean color.RGBA, rms float64) {
	return p.stats(r.X, r.Y, r.W, r.H)
}

func (p *Pixels) heterogeneous(r Rect, opts Options) (color.RGBA, bool) {
	mean, rms := p.stats(r.X, r.Y, r.W, r.H)
	if rms > opts.Threshold {
		return mean, true
	}

	if opts.Detail > 0 && opts.DetailFrac > 0 &&
		p.outlierFraction(r.X, r.Y, r.W, r.H, mean, opts.Detail) > opts.DetailFrac {
		return mean, true
	}

	return mean, false
}

func (p *Pixels) stats(x, y, w, h int) (mean color.RGBA, rms float64) {
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

func (p *Pixels) outlierFraction(x, y, w, h int, mean color.RGBA, dist float64) float64 {
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
