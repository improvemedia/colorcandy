package colorcandy

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"testing"
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

	r, g, b, a := c.RGBA()
	if r != c[0] || g != c[1] || b != c[2] || a != c[3] {
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
	l, a, b, _ := Lab.Convert(ColorFromString("1a3971")).RGBA()
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
