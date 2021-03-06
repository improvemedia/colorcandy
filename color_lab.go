package colorcandy

import (
	"math"
)

func RgbToLab(_r, _g, _b int32) Color {
	r := normalize(float64(_r))
	g := normalize(float64(_g))
	b := normalize(float64(_b))

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

	return Color{
		int32(math.Floor(((116 * fy) - 16) + 0.5)),
		int32(math.Floor(500*(fx-fy) + 0.5)),
		int32(math.Floor(200*(fy-fz) + 0.5)),
		0,
	}
}

func LabMerge(count1 *ColorCount, count2 *ColorCount) (*ColorCount, *ColorCount) {
	// c1 := Lab.Convert(count1.color)
	// c2 := Lab.Convert(count2.color)

	max := count2
	if count1.Percentage >= count2.Percentage {
		max = count1
		max.Total += count2.Total
	} else {
		max.Total += count1.Total
	}

	min := count2
	if count1.Total <= count2.Total {
		min = count1
	}

	return max, min
}

func labChroma(a, b int32) float64 {
	return math.Sqrt(float64((a * a) + (b * b)))
}

func DeltaE(lab_one Color, lab_other Color) float64 {
	l1, a1, b1 := lab_one.RGB()
	l2, a2, b2 := lab_other.RGB()

	c1, c2 := labChroma(a1, b1), labChroma(a2, b2)
	da := a1 - a2
	db := b1 - b2
	dc := c1 - c2
	dh2 := float64(da*da) + float64(db*db) - dc*dc
	if dh2 < 0 {
		return 10000
	}
	pow25_7 := math.Pow(25, 7)
	k_L := 1.0
	k_C := 1.0
	k_H := 1.0
	c1 = math.Sqrt(float64(a1*a1 + b1*b1))
	c2 = math.Sqrt(float64(a2*a2 + b2*b2))
	c_avg := (c1 + c2) / 2
	g := 0.5 * (1 - math.Sqrt(math.Pow(c_avg, 7)/(math.Pow(c_avg, 7)+pow25_7)))
	l1_ := float64(l1)
	a1_ := (1 + g) * float64(a1)
	b1_ := float64(b1)
	l2_ := float64(l2)
	a2_ := (1 + g) * float64(a2)
	b2_ := float64(b2)
	c1_ := math.Sqrt(a1_*a1_ + b1_*b1_)
	c2_ := math.Sqrt(a2_*a2_ + b2_*b2_)
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

	return math.Sqrt(math.Pow(aj, 2) + math.Pow(ak, 2) + math.Pow(al, 2) + r_t*ak*al)
}
