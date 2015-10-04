package colorcandy

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/improvemedia/colorcandy.git/candy"
)

var config Config

func init() {
	r, err := os.Open("./etc/candy.json")
	if err != nil {
		log.Fatal(err)
	}

	d := json.NewDecoder(r)
	d.Decode(&config)
}

func TestColor(t *testing.T) {
	c := Color{255, 255, 255, 0}

	r, g, b := c.RGB()
	if r != c[0] || g != c[1] || b != c[2] {
		t.FailNow()
	}
	if c.Hex() != "ffffff" {
		t.Logf("c.Hex=%s", c.Hex())
		t.FailNow()
	}

	c2 := Color{255, 255, 255, 0}
	if !c.Equal(c2) {
		t.FailNow()
	}

	c3 := ColorFromString("ffffff")
	if c3.Hex() != "ffffff" {
		t.FailNow()
	}
}

func TestRGBToLab(t *testing.T) {
	l, a, b := Lab.Convert(ColorFromString("1a3971")).RGB()
	if l != 24 || a != 5 || b != -36 {
		t.Fatalf("%d,%d,%d", l, a, b)
	}
}

func TestDeltaE(t *testing.T) {
	var deltaE float64
	deltaE = DeltaE(Lab.Convert(ColorFromString("660000")), Lab.Convert(ColorFromString("1a3971")))
	t.Logf("%.2f", deltaE)
	if deltaE != 39.268780242048216 {
		t.FailNow()
	}
	deltaE = DeltaE(Lab.Convert(ColorFromString("cc0000")), Lab.Convert(ColorFromString("1a3971")))
	t.Logf("%.2f", deltaE)
	if deltaE != 47.614294682091405 {
		t.FailNow()
	}
	deltaE = DeltaE(Lab.Convert(ColorFromString("ce454c")), Lab.Convert(ColorFromString("1a3971")))
	t.Logf("%.2f", deltaE)
	if deltaE != 45.56294566011809 {
		t.FailNow()
	}
	deltaE = DeltaE(Lab.Convert(ColorFromString("ea4c88")), Lab.Convert(ColorFromString("1a3971")))
	t.Logf("%.2f", deltaE)
	if deltaE != 47.553035609710335 {
		t.FailNow()
	}
	deltaE = DeltaE(Lab.Convert(ColorFromString("993399")), Lab.Convert(ColorFromString("1a3971")))
	t.Logf("%.2f", deltaE)
	if deltaE != 29.907111520674547 {
		t.FailNow()
	}
	deltaE = DeltaE(Lab.Convert(ColorFromString("663399")), Lab.Convert(ColorFromString("1a3971")))
	t.Logf("%.2f", deltaE)
	if deltaE != 17.41339210323694 {
		t.FailNow()
	}

	deltaE = DeltaE(Lab.Convert2(ColorFromString("af0d09")), Lab.Convert2(ColorFromString("ab1001")))
	t.Logf("%.2f", deltaE)
	if deltaE != 1.5962745605563358 {
		t.FailNow()
	}
}

func TestClosestBaseColor(t *testing.T) {
	c := New(config)
	c1, d1 := c.closestBaseColorTo(ColorFromString("5c2f04"))
	if c1.Hex() != "663300" {
		t.Fatalf("%s != 663300", c1.Hex())
	}
	if d1 != 0 {
		t.Fatalf("d1=%s", d1)
	}
	c2, d2 := c.closestBaseColorTo(ColorFromString("1a3971"))
	if c2.Hex() != "0066cc" {
		t.Fatalf("%s != 0066cc", c2.Hex())
	}
	if math.Floor(d2) != 18 {
		t.Fatalf("d2=%.2f", d2)
	}
}

func TestCompactToCommonColors(t *testing.T) {
	var c Color
	palette := map[Color]*ColorCount{}
	c = ColorFromString("0b100d")
	palette[c] = &ColorCount{c, 61660, 36.27058823529412}
	c = ColorFromString("0f100a")
	palette[c] = &ColorCount{c, 32885, 19.344117647058823}
	c = ColorFromString("0e1117")
	palette[c] = &ColorCount{c, 17194, 10.114117647058823}
	c = ColorFromString("141111")
	palette[c] = &ColorCount{c, 8207, 4.82764705882353}
	c = ColorFromString("352b16")
	palette[c] = &ColorCount{c, 2817, 1.6570588235294117}
	c = ColorFromString("212f2b")
	palette[c] = &ColorCount{c, 3469, 2.0405882352941176}
	c = ColorFromString("470d04")
	palette[c] = &ColorCount{c, 1368, 0.8047058823529412}
	c = ColorFromString("661e0d")
	palette[c] = &ColorCount{c, 1358, 0.7988235294117647}
	c = ColorFromString("5c2f04")
	palette[c] = &ColorCount{c, 1795, 1.0558823529411765}
	c = ColorFromString("422527")
	palette[c] = &ColorCount{c, 1600, 0.9411764705882353}
	c = ColorFromString("754b10")
	palette[c] = &ColorCount{c, 1264, 0.7435294117647059}
	c = ColorFromString("574b2a")
	palette[c] = &ColorCount{c, 1091, 0.6417647058823529}
	c = ColorFromString("242f49")
	palette[c] = &ColorCount{c, 1228, 0.7223529411764706}
	c = ColorFromString("1a3971")
	palette[c] = &ColorCount{c, 2140, 1.2588235294117647}
	c = ColorFromString("135241")
	palette[c] = &ColorCount{c, 1141, 0.6711764705882353}
	c = ColorFromString("3a5555")
	palette[c] = &ColorCount{c, 739, 0.43470588235294116}
	c = ColorFromString("3e4f73")
	palette[c] = &ColorCount{c, 1945, 1.1441176470588235}
	c = ColorFromString("2e6074")
	palette[c] = &ColorCount{c, 823, 0.48411764705882354}
	c = ColorFromString("5b4646")
	palette[c] = &ColorCount{c, 949, 0.558235294117647}
	c = ColorFromString("4f6770")
	palette[c] = &ColorCount{c, 377, 0.22176470588235295}
	c = ColorFromString("931305")
	palette[c] = &ColorCount{c, 2759, 1.6229411764705883}
	c = ColorFromString("af0d09")
	palette[c] = &ColorCount{c, 1074, 0.6317647058823529}
	c = ColorFromString("ab1001")
	palette[c] = &ColorCount{c, 791, 0.4652941176470588}
	c = ColorFromString("80361f")
	palette[c] = &ColorCount{c, 702, 0.41294117647058826}
	c = ColorFromString("a73309")
	palette[c] = &ColorCount{c, 1211, 0.7123529411764706}
	c = ColorFromString("be3004")
	palette[c] = &ColorCount{c, 1095, 0.6441176470588236}
	c = ColorFromString("994607")
	palette[c] = &ColorCount{c, 2026, 1.1917647058823528}
	c = ColorFromString("bc440c")
	palette[c] = &ColorCount{c, 337, 0.19823529411764707}
	c = ColorFromString("907a01")
	palette[c] = &ColorCount{c, 1578, 0.928235294117647}
	c = ColorFromString("956e19")
	palette[c] = &ColorCount{c, 540, 0.3176470588235294}
	c = ColorFromString("a97305")
	palette[c] = &ColorCount{c, 1750, 1.0294117647058822}
	c = ColorFromString("b67c0a")
	palette[c] = &ColorCount{c, 269, 0.15823529411764706}
	c = ColorFromString("c74007")
	palette[c] = &ColorCount{c, 1211, 0.7123529411764706}
	c = ColorFromString("98644c")
	palette[c] = &ColorCount{c, 484, 0.2847058823529412}
	c = ColorFromString("8a7057")
	palette[c] = &ColorCount{c, 1190, 0.7}
	c = ColorFromString("a17a55")
	palette[c] = &ColorCount{c, 611, 0.3594117647058824}
	c = ColorFromString("a3775e")
	palette[c] = &ColorCount{c, 442, 0.26}
	c = ColorFromString("8b6d62")
	palette[c] = &ColorCount{c, 938, 0.5517647058823529}
	c = ColorFromString("9c7e62")
	palette[c] = &ColorCount{c, 812, 0.4776470588235294}
	c = ColorFromString("9e7b6c")
	palette[c] = &ColorCount{c, 1407, 0.8276470588235294}
	c = ColorFromString("2f5894")
	palette[c] = &ColorCount{c, 3026, 1.78}
	c = ColorFromString("4f6382")
	palette[c] = &ColorCount{c, 413, 0.24294117647058824}
	c = ColorFromString("41649e")
	palette[c] = &ColorCount{c, 1284, 0.7552941176470588}

	sample := sampleHistogram()

	compacted := CompactToCommonColors(palette, 2.5)
	for k, v := range compacted {
		if s, ok := sample[k.Hex()]; ok {
			if s.Percentage != v.Percentage {
				t.Fatalf("k: %s %.2f != %.2f", k.Hex(), s.Percentage, v.Percentage)
			}
		}
	}
}

func TestExtractColors(t *testing.T) {
	c := New(config)
	sample := sampleHistogram()
	histogram := map[Color]*ColorCount{}
	for k, v := range sample {
		c := ColorFromString(k)
		histogram[c] = &ColorCount{c, v.Total, v.Percentage}
	}

	var k string
	colors := map[string]*candy.ColorMeta{}

	k = "000000"
	colors[k] = &candy.ColorMeta{
		Colors:       []string{"0b100d", "0f100a", "0e1117", "141111", "352b16", "212f2b", "422527", "5b4646"},
		BaseColor:    "000000",
		SearchFactor: 75,
		Distance:     10000,
	}

	k = "660000"
	colors[k] = &candy.ColorMeta{
		Colors:       []string{"470d04", "661e0d", "931305"},
		BaseColor:    "660000",
		SearchFactor: 3,
		Distance:     0.0,
	}

	k = "663300"
	colors[k] = &candy.ColorMeta{
		Colors:       []string{"5c2f04", "754b10", "574b2a", "80361f"},
		BaseColor:    "663300",
		SearchFactor: 2,
		Distance:     0.0,
	}

	k = "0066cc"
	colors[k] = &candy.ColorMeta{
		Colors:       []string{"242f49", "1a3971", "41649e"},
		BaseColor:    "0066cc",
		SearchFactor: 2,
		Distance:     18.62195249128642,
	}

	k = "304961"
	colors[k] = &candy.ColorMeta{
		Colors:       []string{"3a5555", "3e4f73", "2e6074", "4f6770", "2f5894", "4f6382"},
		BaseColor:    "304961",
		SearchFactor: 4,
		Distance:     0.0,
	}

	k = "cc0000"
	colors[k] = &candy.ColorMeta{
		Colors:       []string{"af0d09", "a73309", "be3004", "bc440c", "c74007"},
		BaseColor:    "cc0000",
		SearchFactor: 3,
		Distance:     0.0,
	}

	k = "996633"
	colors[k] = &candy.ColorMeta{
		Colors:       []string{"994607", "956e19", "a97305", "b67c0a", "98644c", "8a7057", "a17a55", "a3775e", "8b6d62", "9c7e62", "9e7b6c"},
		BaseColor:    "996633",
		SearchFactor: 6,
		Distance:     0.0,
	}

	_colors, _, _ := c.extractColorsFromHistogram(histogram)
	for k, v := range _colors {
		if s, ok := colors[k]; ok {
			t.Logf("key=%s value=%+v sample=%+v", k, v, s)

			check := map[string]bool{}
			for _, v := range v.Colors {
				check[v] = true
			}
			for _, v := range s.Colors {
				if _, ok := check[v]; !ok {
					t.Fatal("Color")
				}
			}

			if s.BaseColor != v.BaseColor {
				t.Fatal("BaseColor")
			}
			if s.Distance != v.Distance {
				t.Fatal("Distance")
			}
			if s.SearchFactor != math.Floor(v.SearchFactor) {
				t.Fatal("SearchFactor")
			}
		} else {
			t.Logf("key %s not found", k)
		}
	}
}

func _TestImages(t *testing.T) {
	c := New(config)
	for i := 0; i < 1; i++ {
		path := "./img/" + strconv.Itoa(i) + ".jpg"
		_, meta, err := c.ExtractColors(path)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range meta {
			fmt.Printf("%s:\n\t%+v\n", k, v)
		}
	}
	t.FailNow()
}

func sampleHistogram() map[string]*ColorCount {
	sample := map[string]*ColorCount{}

	var k string
	var c Color

	k = "0b100d"
	sample[k] = &ColorCount{c, 61660, 36.27058823529412}
	k = "0f100a"
	sample[k] = &ColorCount{c, 32885, 19.344117647058823}
	k = "0e1117"
	sample[k] = &ColorCount{c, 17194, 10.114117647058823}
	k = "141111"
	sample[k] = &ColorCount{c, 8207, 4.82764705882353}
	k = "352b16"
	sample[k] = &ColorCount{c, 2817, 1.6570588235294117}
	k = "212f2b"
	sample[k] = &ColorCount{c, 3469, 2.0405882352941176}
	k = "470d04"
	sample[k] = &ColorCount{c, 1368, 0.8047058823529412}
	k = "661e0d"
	sample[k] = &ColorCount{c, 1358, 0.7988235294117647}
	k = "5c2f04"
	sample[k] = &ColorCount{c, 1795, 1.0558823529411765}
	k = "422527"
	sample[k] = &ColorCount{c, 1600, 0.9411764705882353}
	k = "754b10"
	sample[k] = &ColorCount{c, 1264, 0.7435294117647059}
	k = "574b2a"
	sample[k] = &ColorCount{c, 1091, 0.6417647058823529}
	k = "242f49"
	sample[k] = &ColorCount{c, 1228, 0.7223529411764706}
	k = "1a3971"
	sample[k] = &ColorCount{c, 2140, 1.2588235294117647}
	k = "135241"
	sample[k] = &ColorCount{c, 1141, 0.6711764705882353}
	k = "3a5555"
	sample[k] = &ColorCount{c, 739, 0.43470588235294116}
	k = "3e4f73"
	sample[k] = &ColorCount{c, 1945, 1.1441176470588235}
	k = "2e6074"
	sample[k] = &ColorCount{c, 823, 0.48411764705882354}
	k = "5b4646"
	sample[k] = &ColorCount{c, 949, 0.558235294117647}
	k = "4f6770"
	sample[k] = &ColorCount{c, 377, 0.22176470588235295}
	k = "931305"
	sample[k] = &ColorCount{c, 2759, 1.6229411764705883}
	k = "af0d09"
	sample[k] = &ColorCount{c, 1074, 1.0970588235294116}
	k = "80361f"
	sample[k] = &ColorCount{c, 702, 0.41294117647058826}
	k = "a73309"
	sample[k] = &ColorCount{c, 1211, 0.7123529411764706}
	k = "be3004"
	sample[k] = &ColorCount{c, 1095, 0.6441176470588236}
	k = "994607"
	sample[k] = &ColorCount{c, 2026, 1.1917647058823528}
	k = "bc440c"
	sample[k] = &ColorCount{c, 337, 0.19823529411764707}
	k = "907a01"
	sample[k] = &ColorCount{c, 1578, 0.928235294117647}
	k = "956e19"
	sample[k] = &ColorCount{c, 540, 0.3176470588235294}
	k = "a97305"
	sample[k] = &ColorCount{c, 1750, 1.0294117647058822}
	k = "b67c0a"
	sample[k] = &ColorCount{c, 269, 0.15823529411764706}
	k = "c74007"
	sample[k] = &ColorCount{c, 1211, 0.7123529411764706}
	k = "98644c"
	sample[k] = &ColorCount{c, 484, 0.2847058823529412}
	k = "8a7057"
	sample[k] = &ColorCount{c, 1190, 0.7}
	k = "a17a55"
	sample[k] = &ColorCount{c, 611, 0.3594117647058824}
	k = "a3775e"
	sample[k] = &ColorCount{c, 442, 0.26}
	k = "8b6d62"
	sample[k] = &ColorCount{c, 938, 0.5517647058823529}
	k = "9c7e62"
	sample[k] = &ColorCount{c, 812, 0.4776470588235294}
	k = "9e7b6c"
	sample[k] = &ColorCount{c, 1407, 0.8276470588235294}
	k = "2f5894"
	sample[k] = &ColorCount{c, 3026, 1.78}
	k = "4f6382"
	sample[k] = &ColorCount{c, 413, 0.24294117647058824}
	k = "41649e"
	sample[k] = &ColorCount{c, 1284, 0.7552941176470588}
	return sample
}
