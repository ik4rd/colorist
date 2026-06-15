package palette

import (
	"image"
	"math"
	"sort"
)

const (
	maxSample = 96
	maxCand   = 48
	paletteN  = 4
)

type Accent struct {
	Base   [3]int `json:"base"`
	Deep   [3]int `json:"deep"`
	Bright [3]int `json:"bright"`
	Soft   [3]int `json:"soft"`
}

type Theme struct {
	Palette [][3]int `json:"palette"`
	Accent  *Accent  `json:"accent"`
}

func FromImage(img image.Image) Theme {
	pal := extractPalette(img, paletteN)
	return Theme{Palette: pal, Accent: deriveAccent(pal)}
}

type bucket struct {
	r, g, b, n int
}

type candidate struct {
	rgb [3]int
	n   int
	key int
}

func extractPalette(img image.Image, count int) [][3]int {
	bnds := img.Bounds()
	sw, sh := bnds.Dx(), bnds.Dy()
	if sw <= 0 || sh <= 0 {
		return nil
	}

	scale := math.Min(1, float64(maxSample)/float64(max(sw, sh)))
	tw := max(1, int(math.Round(float64(sw)*scale)))
	th := max(1, int(math.Round(float64(sh)*scale)))

	buckets := make(map[int]*bucket)
	for ty := range th {
		sy := bnds.Min.Y + ty*sh/th

		for tx := range tw {
			sx := bnds.Min.X + tx*sw/tw

			r16, g16, b16, a16 := img.At(sx, sy).RGBA()
			if a16>>8 < 128 {
				continue
			}

			r, g, b := int(r16>>8), int(g16>>8), int(b16>>8)
			key := (r>>3)<<10 | (g>>3)<<5 | (b >> 3)

			e := buckets[key]
			if e == nil {
				e = &bucket{}
				buckets[key] = e
			}

			e.r += r
			e.g += g
			e.b += b
			e.n++
		}
	}

	cands := make([]candidate, 0, len(buckets))

	for key, e := range buckets {
		cands = append(cands, candidate{
			rgb: [3]int{
				int(math.Round(float64(e.r) / float64(e.n))),
				int(math.Round(float64(e.g) / float64(e.n))),
				int(math.Round(float64(e.b) / float64(e.n))),
			},
			n:   e.n,
			key: key,
		})
	}

	if len(cands) == 0 {
		return nil
	}

	sort.Slice(cands, func(i, j int) bool {
		if cands[i].n != cands[j].n {
			return cands[i].n > cands[j].n
		}
		return cands[i].key < cands[j].key
	})

	if len(cands) > maxCand {
		cands = cands[:maxCand]
	}

	chosen := [][3]int{cands[0].rgb}

	for len(chosen) < count && len(chosen) < len(cands) {
		var best [3]int
		bestScore := -1

		for _, c := range cands {
			dmin := math.MaxInt
			for _, s := range chosen {
				if d := dist2(c.rgb, s); d < dmin {
					dmin = d
				}
			}
			if dmin > bestScore {
				bestScore = dmin
				best = c.rgb
			}
		}

		if bestScore == 0 {
			break
		}

		chosen = append(chosen, best)
	}

	for len(chosen) < count {
		chosen = append(chosen, chosen[len(chosen)-1])
	}

	return chosen
}

func deriveAccent(pal [][3]int) *Accent {
	var bh, bs, bl float64
	bestScore := 0.0
	found := false

	for _, c := range pal {
		h, s, l := rgbToHSL(c[0], c[1], c[2])
		if score := s * (1 - math.Abs(l-0.5)); score > bestScore {
			bestScore, bh, bs, bl = score, h, s, l
			found = true
		}
	}

	if !found || bs < 0.12 {
		return nil
	}

	sat := clampF(math.Max(bs, 0.5), 0, 1)

	return &Accent{
		Base:   hslToRGB(bh, sat, clampF(bl, 0.34, 0.44)),
		Deep:   hslToRGB(bh, sat, clampF(bl-0.1, 0.22, 0.4)),
		Bright: hslToRGB(bh, math.Min(1, sat+0.1), 0.58),
		Soft:   hslToRGB(bh, math.Min(0.45, sat*0.4+0.08), 0.92),
	}
}

func dist2(a, b [3]int) int {
	dr, dg, db := a[0]-b[0], a[1]-b[1], a[2]-b[2]
	return dr*dr + dg*dg + db*db
}

func rgbToHSL(ri, gi, bi int) (h, s, l float64) {
	r := float64(ri) / 255
	g := float64(gi) / 255
	b := float64(bi) / 255

	maxv := math.Max(r, math.Max(g, b))
	minv := math.Min(r, math.Min(g, b))
	l = (maxv + minv) / 2

	if maxv == minv {
		return 0, 0, l
	}

	d := maxv - minv
	if l > 0.5 {
		s = d / (2 - maxv - minv)
	} else {
		s = d / (maxv + minv)
	}

	switch maxv {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/d + 2
	default:
		h = (r-g)/d + 4
	}

	return h / 6, s, l
}

func hslToRGB(h, s, l float64) [3]int {
	if s == 0 {
		v := int(math.Round(l * 255))
		return [3]int{v, v, v}
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	return [3]int{
		int(math.Round(hue(p, q, h+1.0/3) * 255)),
		int(math.Round(hue(p, q, h) * 255)),
		int(math.Round(hue(p, q, h-1.0/3) * 255)),
	}
}

func hue(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}

	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 1.0/2:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}

func clampF(v, lo, hi float64) float64 {
	return math.Min(hi, math.Max(lo, v))
}
