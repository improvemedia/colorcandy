package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
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

func ColorFromString(c string) Color {
	var rgb [3]int32
	for i, _ := range rgb {
		b, _ := hex.DecodeString(string(c[i*2]) + string(c[i*2+1]))
		rgb[i] = int32(b[0])
	}
	return Color{rgb[0], rgb[1], rgb[2], 0}
}

func histogram(path string) map[string]*ColorCount {
	var sum float64 = 0.0
	histogram := map[string]*ColorCount{}
	out := convert(path)
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
		histogram[color.Hex()] = &ColorCount{
			int64(count),
			0.0,
		}
	}
	for k, _ := range histogram {
		histogram[k].Percentage = float64(histogram[k].Total) / sum * 100.0
	}

	return histogram
}

func convert(path string) bytes.Buffer {
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
