package main

// TODO(Alexander Yunin): make color type for SRGB
// add func to convert 16bit colors to 8bit
// use new types for something else (for map with pix-count-percent ?)

import (
	"encoding/hex"
	"flag"
	"github.com/gographics/imagick/imagick"
	"log"
	"math"
	"regexp"
)

var infile = flag.String("infile", "img.jpg", "path to image (gif, jpeg, png)")

type color_meta struct {
	id             int
	search_factor  []float64
	distance       float64
	hex            string
	original_color map[string]float64
	hex_of_base    int
}

type cluster_delta struct {
	cluster string
	delta   float64
}

type magix_pixel struct {
	pw    *imagick.PixelWand
	count [2]float64
}

func (c *color_meta) AddSearchFactor(f float64) bool {
	c.search_factor = append(c.search_factor, f)
	return true
}

func (c *color_meta) ReduceSearchFactor() bool {
	acc := 0.0
	for _, v := range c.search_factor {
		acc += v
	}
	c.search_factor[0] = acc
	return true
}

func color(c imagick.Quantum) uint8 {
	if c/255 > 0 {
		c /= 257
	}
	return uint8(c)
}

func (mpix *magix_pixel) Red() uint8 {
	return color(mpix.pw.GetRedQuantum())
}

func (mpix *magix_pixel) Green() uint8 {
	return color(mpix.pw.GetGreenQuantum())
}

func (mpix *magix_pixel) Blue() uint8 {
	return color(mpix.pw.GetBlueQuantum())
}

func (mpix *magix_pixel) GetRgb() [3]uint8 {
	return [3]uint8{mpix.Red(), mpix.Green(), mpix.Blue()}
}

func main() {
	const delta float64 = 2.5

	flag.Parse()
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImage(*infile)
	if err != nil {
		log.Fatal(err)
	}
	mw.QuantizeImage(60, imagick.COLORSPACE_YIQ, 0, true, false)
	_, pixels := mw.GetImageHistogram()

	var sum_of_pixels float64 = 0.0
	delta_count := 0
	mapped_palette := make(map[imagick.PixelWand][2]float64)
	for _, pix := range pixels {
		if pix.IsVerified() == true {
			count := float64(pix.GetColorCount())
			sum_of_pixels += count
			mapped_palette[pix] = [2]float64{count, 0.0}
		}
		defer pix.Destroy()
	}

	for k, v := range mapped_palette {
		mapped_palette[k] = [2]float64{v[0], (v[0] / (sum_of_pixels / 100.0))}
	}

	cp_mapped_palette := mapped_palette
	new_palette := make([]magix_pixel, 0, len(mapped_palette))

	// cant iterate through map values inside a cycle after assigning new values
	// I have to launch new cycle or will get mapped_palette[k] -> nil %)
	for k, v := range mapped_palette {
		common_colors := make([]magix_pixel, 0, len(mapped_palette))
		common_colors = append(common_colors, magix_pixel{&k, v})
		r1, g1, b1 := k.GetRedQuantum(), k.GetGreenQuantum(), k.GetBlueQuantum()
		if r1/255 > 0 {
			r1 /= 257
		}
		if g1/255 > 0 {
			g1 /= 257
		}
		if b1/255 > 0 {
			b1 /= 257
		}
		lab := RgbToLab([3]float64{float64(r1), float64(b1), float64(g1)}) //care, bug right here

		for k2, v2 := range cp_mapped_palette {
			r2, g2, b2 := k2.GetRedQuantum(), k2.GetGreenQuantum(), k2.GetBlueQuantum()
			if r2/255 > 0 {
				r2 /= 257
			}
			if g2/255 > 0 {
				g2 /= 257
			}
			if b2/255 > 0 {
				b2 /= 257
			}
			lab2 := RgbToLab([3]float64{float64(r2), float64(b2), float64(g2)}) //care, bug right here

			if delta > DeltaE(lab, lab2) {
				delta_count += 1 //debug
				if k.GetColorAsString() != k2.GetColorAsString() {
					common_colors = append(common_colors, magix_pixel{&k2, v2})
					common_colors[0].count[1] += v2[1]
				}
			}
		}

		for _, mpix := range common_colors[1:] {
			delete(cp_mapped_palette, *mpix.pw)
		}
		new_palette = append(new_palette, common_colors[0])
	}
	colors_hex := make(map[string]magix_pixel)
	//colors := make(map[int]color_meta)
	base_colors := []string{"660000", "cc0000", "ea4c88", "993399", "663399", "304961", "0066cc", "66cccc", "77cc33", "336600", "cccc33", "ffcc33", "fff533", "ff6600", "c8ad7f", "996633", "663300", "000000", "999999", "cccccc", "ffffff"}
	for _, mpix := range new_palette {
		rgb := mpix.GetRgb()
		cluster_and_delta := ClosestColorToC(rgb)
		hex_color := hex.EncodeToString([]byte{rgb[0], rgb[1], rgb[2]})
		colors_hex["#"+hex_color] = mpix
		var id int
		for i, v := range base_colors { // FIXME:(Alexander Yunin): if defined?(Rails) { SearchColor.find_or_create_by(color: color).id }
			if v == cluster_and_delta.cluster {
				id = i
				break
			}
		}
		log.Println(id)
	}
}
func rad2deg(rad float64) float64 {
	rad = (rad / math.Pi) * 180.0
	return rad
}

func deg2rad(deg float64) float64 {
	deg = (deg / 180.0) * math.Pi
	return deg
}

func LabHue(a, b float64) float64 {
	var ret float64
	if a == 0 && b == 0 {
		ret = 0
	} else {
		ret = float64(int(rad2deg(math.Atan2(b, a))) % 360)
	}
	return ret
}

func LabChroma(a, b float64) float64 {
	chr := math.Sqrt((a * a) + (b * b))
	return chr
}

func DeltaE(lab_one [3]float64, lab_other [3]float64) float64 {
	l1, a1, b1 := lab_one[0], lab_one[1], lab_one[2]
	l2, a2, b2 := lab_other[0], lab_other[1], lab_other[2]
	c1, c2 := LabChroma(a1, b1), LabChroma(a2, b2)
	da := a1 - a2
	db := b1 - b2
	dc := c1 - c2
	dh2 := math.Pow(da, 2) + math.Pow(db, 2) - math.Pow(dc, 2)
	if dh2 < 0 {
		return 10000
	}
	pow25_7 := math.Pow(25, 7)
	k_L := 1.0
	k_C := 1.0
	k_H := 1.0
	c1 = math.Sqrt(math.Pow(a1, 2) + math.Pow(b1, 2))
	c2 = math.Sqrt(math.Pow(a2, 2) + math.Pow(b2, 2))
	c_avg := (c1 + c2) / 2
	g := 0.5 * (1 - math.Sqrt(math.Pow(c_avg, 7)/(math.Pow(c_avg, 7)+pow25_7)))
	l1_ := l1
	a1_ := (1 + g) * a1
	b1_ := b1
	l2_ := l2
	a2_ := (1 + g) * a2
	b2_ := b2
	c1_ := math.Sqrt(math.Pow(a1_, 2) + math.Pow(b1_, 2))
	c2_ := math.Sqrt(math.Pow(a2_, 2) + math.Pow(b2_, 2))
	var h1_, pl, h2_, dh_cond, dh_, dl_, dc_, l__avg, c__avg, h__avg_cond float64
	if a1_ == 0 && b1_ == 0 {
		h1_ = 0.0
	} else {
		if b1_ >= 0 {
			pl = 0.0
		} else {
			pl = 360.0
		}
		h1_ = rad2deg(math.Atan2(b1_, a1_)) + pl
	}
	if a2_ == 0 && b2_ == 0 {
		h2_ = 0
	} else {
		if b2_ >= 0 {
			pl = 0
		} else {
			pl = 360.0
		}
		h2_ = rad2deg(math.Atan2(b2_, a2_)) + pl
	}

	if h2_-h1_ > 180 {
		dh_cond = 1.0
	} else {
		if h2_-h1_ < -180 {
			dh_cond = 2.0
		} else {
			dh_cond = 0
		}
	}

	if dh_cond == 0 {
		dh_ = h2_ - h1_
	} else {
		if dh_cond == 1 {
			dh_ = h2_ - h1_ - 360.0
		} else {
			dh_ = h2_ + 360.0 - h1_
		}
	}

	dl_ = l2_ - l1_
	dc_ = c2_ - c1_
	dc = dc_
	dh_ = 2 * math.Sqrt(c1_*c2_) * math.Sin(deg2rad(dh_/2.0))
	//dh = dh_
	l__avg = math.Floor((l1_ + l2_) / 2) //hack
	c__avg = (c1_ + c2_) / 2
	if c1_*c2_ == 0 {
		h__avg_cond = 3.0
	} else {
		if math.Abs(h2_-h1_) <= 180 {
			h__avg_cond = 0
		} else {
			if h2_+h1_ < 360 {
				h__avg_cond = 1.0
			} else {
				h__avg_cond = 2.0
			}
		}
	}
	var h__avg float64
	if h__avg_cond == 3 {
		h__avg = h1_ + h2_
	} else {
		if h__avg_cond == 0 {
			h__avg = (h1_ + h2_) / 2
		} else {
			if h__avg_cond == 1 {
				h__avg = (h1_+h2_)/2 + 180.0
			} else {
				h__avg = (h1_+h2_)/2 - 180.0
			}
		}
	}
	var ab, s_l, s_c, t, s_h, dtheta, r_c, r_t, aj, ak, al float64
	ab = math.Pow((l__avg - 50.0), 2)
	s_l = 1 + 0.015*ab/math.Sqrt(20.0+ab)
	s_c = 1 + 0.045*c__avg
	t = (1 - 0.17*math.Cos(deg2rad(h__avg-30.0)) + 0.24*math.Cos(deg2rad(2.0*h__avg)) + 0.32*math.Cos(deg2rad(3.0*h__avg+6.0)) - 0.2*math.Cos(deg2rad(4*h__avg-63.0)))
	s_h = 1 + 0.015*c__avg*t
	dtheta = 30.0 * math.Exp(-1*math.Pow(((h__avg-275.0)/25.0), 2))
	r_c = 2.0 * math.Sqrt(math.Pow(c__avg, 7)/(math.Pow(c__avg, 7)+pow25_7))
	r_t = -math.Sin(deg2rad(2.0*dtheta)) * r_c
	aj = dl_ / s_l / k_L
	ak = dc_ / s_c / k_C
	al = dh_ / s_h / k_H
	ret := math.Sqrt(math.Pow(aj, 2) + math.Pow(ak, 2) + math.Pow(al, 2) + r_t*ak*al)
	return ret
}

func Normalize(v float64) float64 {
	var res float64
	v /= 255.0
	if v <= 0.04045 {
		res = v / 12
	} else {
		res = math.Pow(((v + 0.055) / 1.055), 2.4)
	}

	return res
}

func RgbToLabBad(rgb [3]float64) [3]float64 {

	x_d65 := 0.9504
	y_d65 := 1.0
	z_d65 := 1.088

	f_x := Lab(rgb[0] / x_d65)
	f_y := Lab(rgb[1] / y_d65)
	f_z := Lab(rgb[2] / z_d65)

	l := 116*f_y - 16
	a := 500 * (f_x - f_y)
	b := 200 * (f_y - f_z)

	return [3]float64{l, a, b}
}

// Be careful. Actualy main passes to this func R B G, NOT RGB. Copied from colorcake.
func RgbToLab(rgb [3]float64) [3]float64 {
	r, g, b := Normalize(rgb[0]), Normalize(rgb[1]), Normalize(rgb[2])

	x := 0.436052025*r + 0.385081593*g + 0.143087414*b
	y := 0.222491598*r + 0.71688606*g + 0.060621486*b
	z := 0.013929122*r + 0.097097002*g + 0.71418547*b

	xr := x / 0.964221
	yr := y
	zr := z / 0.825211

	eps := 216.0 / 24389
	k := 24389.0 / 27

	var fx, fy, fz float64

	if xr > eps {
		fx = math.Pow(xr, (1.0 / 3))
	} else {
		fx = (k*xr + 16) / 116
	}
	if yr > eps {
		fy = math.Pow(yr, (1.0 / 3))
	} else {
		fy = (k*yr + 16) / 116
	}
	if zr > eps {
		fz = math.Pow(zr, (1.0 / 3))
	} else {
		fz = (k*zr + 16) / 116
	}

	l := ((116 * fy) - 16) //2.55 *
	a := 500 * (fx - fy)
	b = 200 * (fy - fz)

	ret := [3]float64{math.Floor(l + 0.5), math.Floor(a + 0.5), math.Floor(b + 0.5)}
	return ret
}

func Lab(t float64) float64 {
	var l float64
	if t > 0.008856 {
		l = math.Pow(t, (1 / 3.0))
	} else {
		l = 7.787*t + (4 / 29.0)
	}
	return l
}

func RgbFromStr(str string) [3]uint8 {
	re := regexp.MustCompile("..")
	var arr []string
	var rgb [3]uint8
	arr = re.FindAllString(str, -1)
	for i, color := range arr {
		b, _ := hex.DecodeString(color)
		rgb[i] = b[0]
	}

	return rgb
}

// create an interface
func arr_conv(c [3]uint8) [3]float64 {

	return [3]float64{float64(c[0]), float64(c[1]), float64(c[2])}
}

func ClosestColorToC(c [3]uint8) cluster_delta {
	cluster_colors := map[string]string{
		"660000": "660000", "cc0000": "cc0000", "ce454c": "cc0000",
		"ea4c88": "ea4c88", "993399": "993399", "663399": "663399",
		"304961": "304961", "405672": "304961", "0066cc": "0066cc",
		"1a3672": "0066cc", "333399": "0066cc", "0099cc": "0066cc",
		"66cccc": "66cccc", "77cc33": "77cc33", "336600": "336600",
		"cccc33": "cccc33", "999900": "cccc33", "ffcc33": "ffcc33",
		"fff533": "fff533", "efd848": "fff533", "ff6600": "ff6600",
		"c8ad7f": "c8ad7f", "ccad37": "c8ad7f", "e0d3ba": "c8ad7f",
		"996633": "996633", "663300": "663300", "000000": "000000",
		"2e2929": "000000", "999999": "999999", "7e8896": "999999",
		"636363": "999999", "cccccc": "cccccc", "afb5ab": "cccccc",
		"ffffff": "ffffff", "dde2e2": "ffffff", "edefeb": "ffffff",
		/*"ffe6e6": "",*/ "ffe6e6": "ffffff", "d5ccc3": "ffffff",
		"f6fce3": "ffffff", "e1f4fa": "ffffff", "e5e1fa": "ffffff",
		"fbe2f1": "ffffff", "fffae6": "ffffff", "ede7cf": "ffffff",
		"cae0e7": "ffffff", "ede1cf": "ffffff",
		"cad3d5": "ffffff"}

	var closest_color string
	var cluster string
	min_delta := math.MaxFloat64 // pls no buf overflow

	float_c := arr_conv(c)
	lab := RgbToLab(float_c)
	for k, _ := range cluster_colors {
		delta := DeltaE(RgbToLab(arr_conv(RgbFromStr(k))), lab)
		if delta < min_delta {
			min_delta = delta
			closest_color = k
			cluster = cluster_colors[k]
		}
	}
	return cluster_delta{cluster, DeltaE(RgbToLab(arr_conv(RgbFromStr(cluster))), RgbToLab(arr_conv(RgbFromStr(closest_color))))}
}
