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

func TestLoadPalette(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	ppu := New(m)

	expect := []color.RGBA{
		color.RGBA{0x00, 0x00, 0x00, 0xFF}, // 0x1F
		color.RGBA{0x6D, 0x6D, 0x6D, 0xFF}, // 0x00
		color.RGBA{0xB6, 0xB6, 0xB6, 0xFF}, // 0x10
		color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // 0x20
	}

	for n, c := range ppu.bgPalettes[0] {
		if c != expect[n] {
			t.Errorf("unmatch palette: %v %v", c, expect[n])
		}
	}
}

func assertScreen(t *testing.T, ppu *PPU) {
	expect := image.NewRGBA(image.Rect(0, 0, 8, 8))
	pattern.PutImage(expect, 0, 0, ppu.patterns[0]['H'], ppu.bgPalettes[0])

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
				t.Errorf("not black color: %v", screen.At(x, y))
			}
		}
	}
}

func TestPattern(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	m.Write(0x2000, 'H')
	ppu := New(m)
	ppu.SetControl2(0x8)

	assertScreen(t, ppu)
}

func TestScreenAddress(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	m.Write(0x2400, 'H')
	ppu := New(m)
	ppu.SetControl1(1)
	ppu.SetControl2(0x8)

	assertScreen(t, ppu)
}

func TestBGShow(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	m.Write(0x2000, 'H')
	ppu := New(m)
	ppu.SetControl2(0)

	screen := ppu.Render()

	black := color.RGBA{0, 0, 0, 0}
	for y := 0; y < 240; y++ {
		for x := 0; x < 256; x++ {
			if screen.At(x, y) != black {
				t.Errorf("not black color: %v", screen.At(x, y))
			}
		}
	}
}
