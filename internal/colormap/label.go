package colormap

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/ik4rd/colorist/internal/assets"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	labelHeightFraction  = 0.22
	labelMaxFontSize     = 36.0
	labelMinFontSize     = 8.0
	labelPaddingFraction = 0.3
	labelDPI             = 72
	labelLumaThreshold   = 140
	labelNameMinCellPx   = 96
	labelLineSpacing     = 0.2
)

var (
	labelDark  = color.RGBA{R: 40, G: 40, B: 40, A: maxChannel}
	labelLight = color.RGBA{R: 235, G: 235, B: 235, A: maxChannel}

	labelFont = mustParseFont(assets.JetBrainsMono)
)

func mustParseFont(b []byte) *opentype.Font {
	f, err := opentype.Parse(b)
	if err != nil {
		panic("colormap: parse builtin font: " + err.Error())
	}

	return f
}

type labeler struct {
	faces map[int]font.Face
}

func newLabeler() *labeler {
	return &labeler{faces: make(map[int]font.Face)}
}

func (l *labeler) face(size int) font.Face {
	if f, ok := l.faces[size]; ok {
		return f
	}

	f, err := opentype.NewFace(labelFont, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     labelDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil
	}

	l.faces[size] = f

	return f
}

func (l *labeler) draw(dst draw.Image, rect image.Rectangle, hexStr, name string, bg color.RGBA) {
	lines := []string{hexStr}

	bigEnough := name != "" &&
		rect.Dx() >= labelNameMinCellPx &&
		rect.Dy() >= labelNameMinCellPx
	if bigEnough {
		lines = []string{hexStr, name}
	}

	for {
		f, pad, ok := l.fit(rect, lines)
		if ok {
			l.render(dst, rect, lines, f, pad, bg)
			return
		}

		if len(lines) > 1 {
			lines = lines[:1]
			continue
		}

		return
	}
}

func (l *labeler) fit(rect image.Rectangle, lines []string) (font.Face, int, bool) {
	n := len(lines)

	size := float64(rect.Dy()) * labelHeightFraction / lineUnit(n)
	if size > labelMaxFontSize {
		size = labelMaxFontSize
	}
	if size < labelMinFontSize {
		return nil, 0, false
	}

	f := l.face(int(size))
	if f == nil {
		return nil, 0, false
	}

	pad := int(size * labelPaddingFraction)
	avail := rect.Dx() - 2*pad
	if avail <= 0 {
		return nil, 0, false
	}

	if maxW := widestLine(f, lines); maxW > avail {
		size *= float64(avail) / float64(maxW)
		if size < labelMinFontSize {
			return nil, 0, false
		}
		if f = l.face(int(size)); f == nil {
			return nil, 0, false
		}
		pad = int(size * labelPaddingFraction)
	}

	if int(size*lineUnit(n))+2*pad > rect.Dy() {
		return nil, 0, false
	}

	return f, pad, true
}

func (l *labeler) render(dst draw.Image, rect image.Rectangle, lines []string, f font.Face, pad int, bg color.RGBA) {
	src := image.NewUniform(labelColor(bg))
	m := f.Metrics()
	lineH := m.Ascent + m.Descent + fixed.I(int(float64(m.Height.Ceil())*labelLineSpacing))

	baseline := fixed.I(rect.Min.Y+pad) + m.Ascent

	for _, line := range lines {
		drawer := &font.Drawer{
			Dst:  dst,
			Src:  src,
			Face: f,
			Dot: fixed.Point26_6{
				X: fixed.I(rect.Min.X + pad),
				Y: baseline,
			},
		}
		drawer.DrawString(line)
		baseline += lineH
	}
}

func widestLine(f font.Face, lines []string) int {
	maxW := 0
	for _, line := range lines {
		if w := font.MeasureString(f, line).Ceil(); w > maxW {
			maxW = w
		}
	}

	return maxW
}

func lineUnit(n int) float64 {
	return float64(n) + float64(n-1)*labelLineSpacing
}

func labelColor(bg color.RGBA) color.RGBA {
	luma := 0.299*float64(bg.R) + 0.587*float64(bg.G) + 0.114*float64(bg.B)
	if luma > labelLumaThreshold {
		return labelDark
	}

	return labelLight
}
