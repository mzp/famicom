package ppu

import (
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/nesfile"
	"github.com/mzp/famicom/pattern"
)

func TestRenderSize(t *testing.T) {
	m := memory.New()
	ppu := New(m)

	rect := ppu.Render().Bounds()

	if (rect.Min != image.Point{0, 0}) {
		t.Error()
	}

	if (rect.Max != image.Point{256, 240}) {
		t.Error()
	}
}

func load(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	rom := nesfile.Load(data)
	return rom.Character
}

func TestPattern(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello.nes"))
	m.Write(0x2000, 'H')
	m.Write(0x2001, 0)
	ppu := New(m)

	expect := image.NewRGBA(image.Rect(0, 0, 8, 8))

	palette := []color.Color{
		color.RGBA{0, 0, 0, 0xFF},
		color.RGBA{0xA0, 0xA0, 0xA0, 0xFF},
		color.RGBA{0, 0xFF, 0, 0xFF},
		color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
	}

	pattern.PutImage(expect, 0, 0, ppu.patterns[0]['H'], palette)

	screen := ppu.Render()

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if screen.At(x, y) != expect.At(x, y) {
				t.Errorf("unmatch color: %v %v", screen.At(x, y), expect.At(x, y))
			}
		}
	}

	black := color.RGBA{0, 0, 0, 0xFF}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if screen.At(x+8, y) != black {
				t.Errorf("unmatch color: %v", screen.At(x, y))
			}
		}
	}
}
