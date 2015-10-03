package colorcandy

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const CompactDelta float64 = 2.5

func CompactToCommonColors(_original map[Color]*ColorCount) map[Color]*ColorCount {
	result := map[Color]*ColorCount{}
	_copy := _original
	deltaCount := 0

	for color1, v1 := range _original {
		commonColors := []*ColorCount{v1}

		for color2, v2 := range _copy {
			if CompactDelta > DeltaE(Lab.Convert(color1), Lab.Convert(color2)) {
				deltaCount += 1 //debug
				if !color1.Equal(color2) {
					commonColors = append(commonColors, v2)
					commonColors[0].Percentage += v2.Percentage
				}
			}
		}

		for _, count := range commonColors[1:] {
			delete(_copy, count.color)
		}

		result[commonColors[0].color] = commonColors[0]
	}

	return result
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
			uint(count),
			0.0,
		}
	}
	for k, _ := range histogram {
		histogram[k].Percentage = float64(histogram[k].Total) / sum * 100.0
	}
	return histogram
}

func _convert(path string) bytes.Buffer {
	cmd := exec.Command("convert", "+dither", "-colors", "60", "-quantize", "YIQ", "-format", "%c", path, "histogram:info:-")
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
