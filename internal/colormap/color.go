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
