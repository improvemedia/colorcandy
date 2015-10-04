package colorcandy

import (
	"bytes"
	_ "fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func ImageHistogram_Cmd(path string) map[Color]*ColorCount {
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
	cmd := exec.Command("convert", "-dither", "Riemersma", "-colors", "60", "-quantize", "YIQ", "-depth", "0", "-format", "%c", path, "histogram:info:-")
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
