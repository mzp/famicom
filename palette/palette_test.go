package palette

import (
	"image/color"
	"testing"
)

func TestRead(t *testing.T) {
	palette := Read([]byte{
		0x0f, 0x11, 0x21, 0x31,
	})

	if len(palette) != 4 {
		t.Errorf("cannot read 4 color: %d", len(palette))
	}

	expect := []color.RGBA{
		color.RGBA{0, 0, 0, 0xFF},
		color.RGBA{0, 0x6D, 0xDA, 0xFF},
		color.RGBA{0x6D, 0xB6, 0xFF, 0xFF},
		color.RGBA{0xB6, 0xDA, 0xFF, 0xFF},
	}

	for n, c := range palette {
		if expect[n] != c {
			t.Errorf("unmatch: %v %v", expect[n], c)
		}
	}
}
