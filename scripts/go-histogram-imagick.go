package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/gographics/imagick/imagick"
	"log"
	"os"
)

var mw *imagick.MagickWand

func main() {
	imagick.Initialize()

	mw = imagick.NewMagickWand()
	defer mw.Destroy()

	h := histogram(os.Args[0])
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(h)
}

type Color [4]int32

func (c Color) Hex() string {
	return hex.EncodeToString([]byte{byte(c[0]), byte(c[1]), byte(c[2])})
}

type ColorCount struct {
	Total      int64
	Percentage float64
}

func histogram(path string) map[string]*ColorCount {
	var sum float64 = 0.0
	histogram := map[string]*ColorCount{}

	err := mw.ReadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	mw.QuantizeImage(60, imagick.COLORSPACE_YIQ, 0, true, false)
	_, pixels := mw.GetImageHistogram()

	for _, pix := range pixels {
		if pix.IsVerified() == true {
			count := int64(pix.GetColorCount())
			sum += float64(count)
			r := pix.GetRedQuantum()
			g := pix.GetGreenQuantum()
			b := pix.GetBlueQuantum()

			color := Color{int32(r), int32(g), int32(b), 0}
			histogram[color.Hex()] = &ColorCount{
				Total:      count,
				Percentage: 0.0,
			}
		}
	}

	for k, v := range histogram {
		histogram[k].Percentage = float64(v.Total) / sum * 100.0
	}

	return histogram
}
