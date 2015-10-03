package colorcandy

import (
	"encoding/hex"
)

type Model struct {
	converter func(Color) Color
}

func (m Model) Convert(c Color) Color {
	return m.converter(c)
}

var (
	Lab Model = Model{RgbToLab}
)

type Color [4]int32

type ColorCount struct {
	color      Color
	Total      int64
	Percentage float64
}

func (c Color) RGBA() (int32, int32, int32, int32) {
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

func NewColor(r int32, g int32, b int32) Color {
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
	var rgb [3]int32
	for i, _ := range rgb {
		b, _ := hex.DecodeString(string(c[i*2]) + string(c[i*2+1]))
		rgb[i] = int32(b[0])
	}
	return Color{rgb[0], rgb[1], rgb[2], 0}
}
