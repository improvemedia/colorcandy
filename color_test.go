package colorcandy

import (
	"testing"
)

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
