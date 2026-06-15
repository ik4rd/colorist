package colormap

import (
	"fmt"
	"image/color"
	"math"
)

const maxChannel = 255

func clamp8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > maxChannel {
		return maxChannel
	}
	return uint8(math.Round(v))
}

func hex(c color.RGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

func rgbStr(c color.RGBA) string {
	return fmt.Sprintf("rgb(%d,%d,%d)", c.R, c.G, c.B)
}

func cmykStr(c color.RGBA) string {
	r := float64(c.R) / maxChannel
	g := float64(c.G) / maxChannel
	b := float64(c.B) / maxChannel

	k := 1 - math.Max(r, math.Max(g, b))
	if k >= 1 {
		return "cmyk(0,0,0,100)"
	}

	cc := (1 - r - k) / (1 - k)
	mm := (1 - g - k) / (1 - k)
	yy := (1 - b - k) / (1 - k)

	return fmt.Sprintf("cmyk(%d,%d,%d,%d)",
		pct(cc), pct(mm), pct(yy), pct(k))
}

func pct(v float64) int {
	return int(math.Round(v * 100))
}

func toRGB(c color.RGBA) string {
	return fmt.Sprintf("rgb(%d,%d,%d)", c.R, c.G, c.B)
}

func toCMYK(c color.RGBA) string {
	if c.R == 0 && c.G == 0 && c.B == 0 {
		return "CMYK(0,0,0,100)"
	}

	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	k := 1 - math.Max(math.Max(r, g), b)
	c_val := (1 - r - k) / (1 - k)
	m := (1 - g - k) / (1 - k)
	y := (1 - b - k) / (1 - k)

	return fmt.Sprintf("CMYK(%d,%d,%d,%d)",
		int(math.Round(c_val*100)),
		int(math.Round(m*100)),
		int(math.Round(y*100)),
		int(math.Round(k*100)))
}
