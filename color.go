package colorcandy

import (
	"encoding/hex"
	"image/color"
)

var (
	Lab color.Model = color.ModelFunc(RgbToLab)
)

type Color [4]uint32

type ColorCount struct {
	color      Color
	Total      uint
	Percentage float64
}

func (c Color) RGBA() (uint32, uint32, uint32, uint32) {
	return c[0], c[1], c[2], c[3]
}

func (c Color) Equal(other Color) bool {
	r1, g1, b1, _ := c.RGBA()
	r2, g2, b2, _ := other.RGBA()

	return r1 == r2 && g1 == g2 && b1 == b2
}

func (c Color) Hex() string {
	return hex.EncodeToString([]byte{byte(c[0]), byte(c[1]), byte(c[2])})
}

func NewColor(r uint32, g uint32, b uint32) Color {
	if r/255 > 0 {
		r /= 257
	}
	if g/255 > 0 {
		g /= 257
	}
	if b/255 > 0 {
		b /= 257
	}
	return Color{r, g, b, 0}
}

func ColorFromString(c string) Color {
	var rgb [3]uint32
	for i, _ := range rgb {
		b, _ := hex.DecodeString(string(c[i*2]) + string(c[i*2+1]))
		rgb[i] = uint32(b[0])
	}
	return Color{rgb[0], rgb[1], rgb[2], 0}
}
