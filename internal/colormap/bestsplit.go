package colormap

import "math"

type Axis int

const (
	AxisVertical Axis = iota
	AxisHorizontal
)

type moments struct {
	sr, sg, sb    float64
	sr2, sg2, sb2 float64
}

func (m *moments) addPixel(r, g, b float64) {
	m.sr += r
	m.sg += g
	m.sb += b
	m.sr2 += r * r
	m.sg2 += g * g
	m.sb2 += b * b
}

func (m *moments) add(o moments) {
	m.sr += o.sr
	m.sg += o.sg
	m.sb += o.sb
	m.sr2 += o.sr2
	m.sg2 += o.sg2
	m.sb2 += o.sb2
}

func (m moments) sub(o moments) moments {
	return moments{
		sr:  m.sr - o.sr,
		sg:  m.sg - o.sg,
		sb:  m.sb - o.sb,
		sr2: m.sr2 - o.sr2,
		sg2: m.sg2 - o.sg2,
		sb2: m.sb2 - o.sb2,
	}
}

func (m moments) sse(n float64) float64 {
	if n <= 0 {
		return 0
	}
	return (m.sr2 - m.sr*m.sr/n) +
		(m.sg2 - m.sg*m.sg/n) +
		(m.sb2 - m.sb*m.sb/n)
}

func (p *Pixels) BestSplit(r Rect, axis Axis, minSize int) (at int, cost float64, ok bool) {
	if minSize < 1 {
		minSize = 1
	}

	lines, cross := r.W, r.H
	if axis == AxisHorizontal {
		lines, cross = r.H, r.W
	}
	if lines < 2*minSize {
		return 0, 0, false
	}

	acc := make([]moments, lines)
	if axis == AxisVertical {
		for j := r.Y; j < r.Y+r.H; j++ {
			row := j * p.w
			for i := range lines {
				k := row + r.X + i
				acc[i].addPixel(p.r[k], p.g[k], p.b[k])
			}
		}
	} else {
		for j := range lines {
			row := (r.Y + j) * p.w
			for i := r.X; i < r.X+r.W; i++ {
				k := row + i
				acc[j].addPixel(p.r[k], p.g[k], p.b[k])
			}
		}
	}

	for i := 1; i < lines; i++ {
		acc[i].add(acc[i-1])
	}
	total := acc[lines-1]

	best := math.Inf(1)
	bestC := -1

	for c := minSize; c <= lines-minSize; c++ {
		left := acc[c-1]
		right := total.sub(left)
		s := left.sse(float64(c*cross)) + right.sse(float64((lines-c)*cross))
		if s < best {
			best = s
			bestC = c
		}
	}

	if bestC < 0 {
		return 0, 0, false
	}

	if axis == AxisVertical {
		return r.X + bestC, best, true
	}

	return r.Y + bestC, best, true
}
