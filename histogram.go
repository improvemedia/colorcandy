package colorcandy

import (
	"bytes"
	_ "fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func CompactToCommonColors(_original map[Color]*ColorCount, delta float64) map[Color]*ColorCount {
	_copy := map[Color]struct{}{}
	for k, _ := range _original {
		_copy[k] = struct{}{}
	}

	for _, v1 := range _original {
		toRemove := []Color{}
		commonColors := map[Color]*ColorCount{}

		for k2, _ := range _copy {
			v2 := _original[k2]
			d := DeltaE(Lab.Convert2(v1.color), Lab.Convert2(v2.color))
			if delta > d {
				if v1.color.Equal(v2.color) {
					if _, ok := commonColors[v1.color]; !ok {
						commonColors[v1.color] = v1
					}
				} else {
					if _, ok := commonColors[v2.color]; !ok {
						commonColors[v2.color] = v2
					}

					if _, ok := commonColors[v1.color]; ok {
						v2.Percentage += v1.Percentage
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

func ImageHistogram(path string) map[Color]*ColorCount {
	var sum float64 = 0.0
	histogram := map[Color]*ColorCount{}
	out := _convert(path)
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line == "" {
			break
		}
		split1 := strings.Split(line, "#")
		a, b := split1[0], split1[1]
		split2 := strings.Split(a, ": ")
		count, _ := strconv.Atoi(strings.TrimLeft(split2[0], " "))
		color := ColorFromString(b[0:6])
		sum += float64(count)
		histogram[color] = &ColorCount{
			color,
			int64(count),
			0.0,
		}
	}
	for k, _ := range histogram {
		histogram[k].Percentage = float64(histogram[k].Total) / sum * 100.0
	}

	return histogram
}

func _convert(path string) bytes.Buffer {
	cmd := exec.Command("convert", "+dither", "-colors", "60", "-quantize", "YIQ", "-depth", "0", "-format", "%c", path, "histogram:info:-")
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Run()
	if err != nil {
		log.Printf("%s: %s", err, errout)
		log.Fatal()
	}
	return out
}
