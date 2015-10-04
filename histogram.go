package colorcandy

import (
	_ "fmt"
)

func CompactToCommonColors(_original map[Color]*ColorCount, delta float64) map[Color]*ColorCount {
	_copy := map[Color]struct{}{}
	for k, _ := range _original {
		_copy[k] = struct{}{}
	}

	for _, v1 := range _original {
		toRemove := []Color{}
		commonColors := map[Color]*ColorCount{
			v1.color: v1,
		}

		for k2, _ := range _copy {
			v2 := _original[k2]
			d := DeltaE(Lab.Convert2(v1.color), Lab.Convert2(v2.color))
			if delta > d {
				if !v1.color.Equal(v2.color) {
					if _, ok := commonColors[v2.color]; !ok {
						commonColors[v2.color] = v2
					}

					if _, ok := commonColors[v1.color]; ok {
						v1.Percentage += v2.Percentage
						toRemove = append(toRemove, v2.color)
					} else {
						commonColors[v1.color] = v1
					}
				}
			}
		}
		for _, k := range toRemove {
			delete(_copy, k)
		}
	}

	return _original
}
