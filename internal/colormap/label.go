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

func (l *labeler) draw(dst draw.Image, rect image.Rectangle, text string, bg color.RGBA) {
	size := float64(rect.Dy()) * labelHeightFraction

	if size > labelMaxFontSize {
		size = labelMaxFontSize
	}
	if size < labelMinFontSize {
		return
	}

	f := l.face(int(size))
	if f == nil {
		return
	}

	pad := int(size * labelPaddingFraction)
	avail := rect.Dx() - 2*pad
	if avail <= 0 {
		return
	}

	if textW := font.MeasureString(f, text).Ceil(); textW > avail {
		size *= float64(avail) / float64(textW)

		if size < labelMinFontSize {
			return
		}
		if f = l.face(int(size)); f == nil {
			return
		}

		pad = int(size * labelPaddingFraction)
	}

	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(labelColor(bg)),
		Face: f,
		Dot: fixed.Point26_6{
			X: fixed.I(rect.Min.X + pad),
			Y: fixed.I(rect.Min.Y+pad) + f.Metrics().Ascent,
		},
	}

	drawer.DrawString(text)
}

func labelColor(bg color.RGBA) color.RGBA {
	luma := 0.299*float64(bg.R) + 0.587*float64(bg.G) + 0.114*float64(bg.B)
	if luma > labelLumaThreshold {
		return labelDark
	}

	return labelLight
}
