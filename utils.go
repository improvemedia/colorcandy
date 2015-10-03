package colorcandy

import (
	"math"
)

func rad2deg(rad float64) float64 {
	return (rad / math.Pi) * 180.0
}

func deg2rad(deg float64) float64 {
	return (deg / 180.0) * math.Pi
}

func normalize(v float64) (ret float64) {
	v /= 255.0
	if v <= 0.04045 {
		ret = v / 12
	} else {
		ret = math.Pow(((v + 0.055) / 1.055), 2.4)
	}
	return
}
