//oh uh
package main

import (
	"flag"
	"github.com/gographics/imagick/imagick"
	"log"
	"math"
)

var infile = flag.String("infile", "img.jpg", "path to image (gif, jpeg, png)")

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
	//log.Println("Histogram R : G : B : COUNT")
	for _, pix := range pixels {
		//log.Println(pix.GetRedQuantum(), ":", pix.GetGreenQuantum(), ":", pix.GetBlueQuantum(), ":", pix.GetColorCount())
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

	//log.Println("After Math")
	// cant iterate through map values inside a cycle after assigning new values
	// I have to launch new cycle or will get mapped_palette[k] -> nil %)
	for k, _ := range mapped_palette {
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

		//log.Println(k.GetRedQuantum(), "->", v, ":", lab)

		for k2, v := range mapped_palette {
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
			log.Println("DELTA FOR", r2, g2, b2, "->", v)
			if delta > DeltaE(lab, lab2) {
				delta_count += 1
			}

		}
	}
	log.Println("delta calc count:", delta_count)

}
func rad2deg(rad float64) float64 {
	//log.Println("rad2deg IN", ":", rad)
	rad = (rad / math.Pi) * 180.0
	//log.Println("rad2deg OUT", ":", rad)
	return rad
}

func deg2rad(deg float64) float64 {
	//log.Println("deg2rad IN", ":", deg)
	deg = (deg / 180.0) * math.Pi
	//log.Println("deg2rad OUT", ":", deg)
	return deg
}

func LabHue(a, b float64) float64 {
	//log.Println("LAbHue IN", ":", a, b)
	var ret float64
	if a == 0 && b == 0 {
		ret = 0
	} else {
		ret = float64(int(rad2deg(math.Atan2(b, a))) % 360)
	}
	//log.Println("LAbHue OUT", ":", ret)
	return ret
}

func LabChroma(a, b float64) float64 {
	//log.Println("LAbChroma IN", ":", a, b)
	chr := math.Sqrt((a * a) + (b * b))
	//log.Println("LAbChroma OUT", ":", chr)
	return chr
}

func DeltaE(lab_one [3]float64, lab_other [3]float64) float64 {
	log.Println("!!!----DELTA_E START----!!!")
	log.Println(lab_one, lab_other)
	log.Println("---------------------------")
	l1, a1, b1 := lab_one[0], lab_one[1], lab_one[2]
	l2, a2, b2 := lab_other[0], lab_other[1], lab_other[2]
	c1, c2 := LabChroma(a1, b1), LabChroma(a2, b2)
	log.Println("c1, c2 | 126:", c1, c2)
	//h1 := LabHue(a1, b1)
	//dl := l2 - l1
	da := a1 - a2
	log.Println("da | 130:", da)
	db := b1 - b2
	log.Println("db | 132:", db)
	dc := c1 - c2
	log.Println("dc | 134:", dc)
	dh2 := math.Pow(da, 2) + math.Pow(db, 2) - math.Pow(dc, 2)
	log.Println("dh2 | 136:", dh2)
	if dh2 < 0 {
		log.Println("!!!----DELTA_E END CAUSE OF 138----!!!")
		return 10000
	}
	//dh := math.Sqrt(dh2)
	//--- case meth
	pow25_7 := math.Pow(25, 7)
	k_L := 1.0
	k_C := 1.0
	k_H := 1.0
	c1 = math.Sqrt(math.Pow(a1, 2) + math.Pow(b1, 2))
	log.Println("c1 | 148:", c1)
	c2 = math.Sqrt(math.Pow(a2, 2) + math.Pow(b2, 2))
	log.Println("c2 | 150:", c2)
	c_avg := (c1 + c2) / 2
	g := 0.5 * (1 - math.Sqrt(math.Pow(c_avg, 7)/(math.Pow(c_avg, 7)+pow25_7)))
	log.Println("g | 153:", g)
	l1_ := l1
	a1_ := (1 + g) * a1
	b1_ := b1
	l2_ := l2
	a2_ := (1 + g) * a2
	b2_ := b2
	c1_ := math.Sqrt(math.Pow(a1_, 2) + math.Pow(b1_, 2))
	log.Println("c1_ | 161:", c1_)
	c2_ := math.Sqrt(math.Pow(a2_, 2) + math.Pow(b2_, 2))
	log.Println("c2 | 163:", c2_)
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
	l__avg = (l1_ + l2_) / 2
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
	log.Println("-------------")
	log.Println("h1_", "pl", "h2_", "dh_cond", "dh_", "dl_", "dc_", "l__avg", "c__avg", "h__avg_cond")
	log.Println(h1_, pl, h2_, dh_cond, dh_, dl_, dc_, l__avg, c__avg, h__avg_cond)
	log.Println("-------------")
	ab = math.Pow((l__avg - 50.0), 2)
	s_l = 1 + 0.015*ab/math.Sqrt(20.0+ab)
	s_c = 1 + 0.045*c__avg
	t = (1 - 0.17*math.Cos(deg2rad(h__avg-30.0)) + 0.24*math.Cos(deg2rad(2.0*h__avg)) + 0.32*math.Cos(deg2rad(3.0*h__avg+6.0)) - 0.2*math.Cos(deg2rad(4*h__avg-63.0)))
	s_h = 1 + 0.015*c__avg*t
	dtheta = 30.0 * math.Pow(math.Exp(-1*((h__avg-275.0)/25.0)), 2)
	r_c = 2.0 * math.Sqrt(math.Pow(c__avg, 7)/(math.Pow(c__avg, 7)+pow25_7))
	r_t = -math.Sin(deg2rad(2.0*dtheta)) * r_c
	aj = dl_ / s_l / k_L
	ak = dc_ / s_c / k_C
	al = dh_ / s_h / k_H
	log.Println("-------------")
	log.Println("ab", "s_l", "s_c", "t", "s_h", "dtheta", "r_c", "r_t", "aj", "ak", "al")
	log.Println(ab, s_l, s_c, t, s_h, dtheta, r_c, r_t, aj, ak, al)
	log.Println("-------------")
	ret := math.Sqrt(math.Pow(aj, 2) + math.Pow(ak, 2) + math.Pow(al, 2) + r_t*ak*al)
	log.Println("ret:", ret)
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
	log.Println("normalized IN:", r, g, b)

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
	log.Println("normalized:", ret)
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
