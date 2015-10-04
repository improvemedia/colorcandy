package colorcandy

import (
	"github.com/gographics/imagick/imagick"
	"log"
)

//var mw *imagick.MagickWand

func init() {
	imagick.Initialize()
}

func ImageHistogram_Imagick(path string) map[Color]*ColorCount {
	var sum float64 = 0.0
	histogram := map[Color]*ColorCount{}

	mw := imagick.NewMagickWand()
	//defer mw.Destroy()

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

			color := NewColor(int32(r), int32(g), int32(b))
			histogram[color] = &ColorCount{
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
