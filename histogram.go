package colorcandy

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func ImageHistogram(path string) map[Color]*ColorCount {
	out := _convert(path)

	var sum_of_pixels float64 = 0.0

	mapped_palette := map[Color]*ColorCount{}

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
		sum_of_pixels += float64(count)
		mapped_palette[color] = &ColorCount{
			color,
			uint(count),
			0.0,
		}
	}
	for k, _ := range mapped_palette {
		mapped_palette[k].Percentage = float64(mapped_palette[k].Total) / (sum_of_pixels / 100.0)
	}
	return mapped_palette
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
