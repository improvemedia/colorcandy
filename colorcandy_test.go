package colorcandy

import (
	"encoding/json"
	"fmt"
	"log"
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

func TestImages(t *testing.T) {
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
