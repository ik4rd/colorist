package colornames

import (
	_ "embed"
	"image/color"
	"math"
	"strings"
	"sync"
)

//go:embed colornames.csv
var raw string

type entry struct {
	name    string
	l, a, b float64
}

var entries []entry

var (
	cacheMu sync.RWMutex
	cache   = make(map[color.RGBA]string)
)

func init() {
	for line := range strings.SplitSeq(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}

		before, after, ok0 := strings.Cut(line, ",")
		if !ok0 {
			continue
		}

		name, hex := before, after
		c, ok := parseHex(hex)
		if !ok {
			continue
		}

		l, a, b := rgbToLab(c)
		entries = append(entries, entry{name: name, l: l, a: a, b: b})
	}
}

func Nearest(c color.RGBA) string {
	c.A = 0xFF

	cacheMu.RLock()
	if name, ok := cache[c]; ok {
		cacheMu.RUnlock()
		return name
	}
	cacheMu.RUnlock()

	l, a, b := rgbToLab(c)

	best := ""
	bestDist := -1.0

	for _, e := range entries {
		dl, da, db := l-e.l, a-e.a, b-e.b
		d := dl*dl + da*da + db*db
		if bestDist < 0 || d < bestDist {
			bestDist = d
			best = e.name
		}
	}

	cacheMu.Lock()
	cache[c] = best
	cacheMu.Unlock()

	return best
}

func parseHex(s string) (color.RGBA, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return color.RGBA{}, false
	}

	r, ok1 := hexByte(s[0], s[1])
	g, ok2 := hexByte(s[2], s[3])
	b, ok3 := hexByte(s[4], s[5])
	if !ok1 || !ok2 || !ok3 {
		return color.RGBA{}, false
	}

	return color.RGBA{R: r, G: g, B: b, A: 0xFF}, true
}

func hexByte(hi, lo byte) (uint8, bool) {
	h, ok1 := hexNibble(hi)
	l, ok2 := hexNibble(lo)
	if !ok1 || !ok2 {
		return 0, false
	}

	return h<<4 | l, true
}

func hexNibble(c byte) (uint8, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

func rgbToLab(c color.RGBA) (l, a, b float64) {
	x, y, z := rgbToXYZ(c)

	const xn, yn, zn = 0.95047, 1.0, 1.08883
	fx := labF(x / xn)
	fy := labF(y / yn)
	fz := labF(z / zn)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)

	return l, a, b
}

func rgbToXYZ(c color.RGBA) (x, y, z float64) {
	r := srgbToLinear(float64(c.R) / 255.0)
	g := srgbToLinear(float64(c.G) / 255.0)
	b := srgbToLinear(float64(c.B) / 255.0)

	x = r*0.4124 + g*0.3576 + b*0.1805
	y = r*0.2126 + g*0.7152 + b*0.0722
	z = r*0.0193 + g*0.1192 + b*0.9505

	return x, y, z
}

func srgbToLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}

	return math.Pow((v+0.055)/1.055, 2.4)
}

func labF(t float64) float64 {
	const e = 216.0 / 24389.0
	const k = 24389.0 / 27.0
	if t > e {
		return math.Cbrt(t)
	}

	return (k*t + 16) / 116
}
