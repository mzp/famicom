package ppu

import (
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mzp/famicom/bits"
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

func TestLoadBgPalette(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	ppu := New(m)

	expect := []color.RGBA{
		color.RGBA{0x00, 0x00, 0x00, 0xFF}, // 0x1F
		color.RGBA{0x6D, 0x6D, 0x6D, 0xFF}, // 0x00
		color.RGBA{0xB6, 0xB6, 0xB6, 0xFF}, // 0x10
		color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // 0x20
	}

	for n, c := range ppu.bg.palettes[0] {
		if c != expect[n] {
			t.Errorf("unmatch palette: %v %v", c, expect[n])
		}
	}
}

func at(img image.Image, x, y int) color.Color {
	return img.At(img.Bounds().Min.X+x, img.Bounds().Min.Y+y)
}

func assertScreen(t *testing.T, ppu *PPU, c byte, dx, dy int, palette []color.Color) {
	expect := image.NewRGBA(image.Rect(0, 0, 8, 8))
	ppu.patterns[0][c].Put(expect, image.Point{}, palette, pattern.Option{
		BackdropColor: palette[0],
	})

	screen := ppu.Render()

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if at(screen, x+dx, y+dy) != expect.At(x, y) {
				t.Errorf("unmatch color: %v %v", at(screen, x+dx, y+dy), expect.At(x, y))
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

	assertScreen(t, ppu, 'H', 0, 0, ppu.bg.palettes[0])
	assertScreen(t, ppu, ' ', 8, 0, ppu.bg.palettes[0])
}

func TestScreenAddress(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	m.Write(0x2400, 'H')
	ppu := New(m)
	ppu.SetControl1(1)
	ppu.SetControl2(0x8)
	ppu.SetVerticalMirror(true)

	assertScreen(t, ppu, 'H', 0, 0, ppu.bg.palettes[0])
	assertScreen(t, ppu, ' ', 8, 0, ppu.bg.palettes[0])
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

func TestVRAWWriteX(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	ppu := New(m)
	ppu.SetControl2(0x8)

	ppu.SetAddress(0x20)
	ppu.SetAddress(0x0)

	ppu.WriteVRAM('H')
	ppu.WriteVRAM('E')

	assertScreen(t, ppu, 'H', 0, 0, ppu.bg.palettes[0])
	assertScreen(t, ppu, 'E', 8, 0, ppu.bg.palettes[0])
}

func TestVRAMReadWithoutBuffer(t *testing.T) {
	m := memory.New()
	m.Write(0x3FFE, 'H')
	m.Write(0x3FFF, 'E')
	ppu := New(m)

	ppu.SetAddress(0x3F)
	ppu.SetAddress(0xFE)
	if ppu.ReadVRAM() != 'H' {
		t.Error()
	}
	if ppu.ReadVRAM() != 'E' {
		t.Error()
	}
}

func TestVRAMReadWithinBuffer(t *testing.T) {
	m := memory.New()
	m.Write(0x2000, 'H')
	m.Write(0x2001, 'E')
	ppu := New(m)

	ppu.SetAddress(0x20)
	ppu.SetAddress(0x0)

	ppu.ReadVRAM()
	if ppu.ReadVRAM() != 'H' {
		t.Error()
	}
	if ppu.ReadVRAM() != 'E' {
		t.Error()
	}
}

func TestVRAWWriteY(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	ppu := New(m)
	ppu.SetControl1(0x4)
	ppu.SetControl2(0x8)

	ppu.SetAddress(0x20)
	ppu.SetAddress(0x0)

	ppu.WriteVRAM('H')
	ppu.WriteVRAM('E')

	assertScreen(t, ppu, 'H', 0, 0, ppu.bg.palettes[0])
	assertScreen(t, ppu, 'E', 0, 8, ppu.bg.palettes[0])
}

func TestPatternSelector(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	m.Write(0x2000, 'H')
	ppu := New(m)
	ppu.SetControl1(0x10)
	ppu.SetControl2(0x8)

	assertScreen(t, ppu, ' ', 0, 0, ppu.bg.palettes[0])
}

func TestLoadSpritePalette(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	ppu := New(m)

	expect := []color.RGBA{
		color.RGBA{0x00, 0x00, 0x00, 0xFF}, // 0x1F
		color.RGBA{0xB6, 0x00, 0x6D, 0xFF}, // 0x05
		color.RGBA{0xFF, 0x00, 0x91, 0xFF}, // 0x15
		color.RGBA{0xFF, 0x6D, 0xFF, 0xFF}, // 0x25
	}

	if len(ppu.sprite.palettes[0]) != 4 {
		t.Error()
	}

	for n, c := range ppu.sprite.palettes[0] {
		if c != expect[n] {
			t.Errorf("unmatch palette: %v %v", c, expect[n])
		}
	}
}

func TestSprite(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))

	ppu := New(m)
	ppu.SetControl2(0x10)
	ppu.SetSpriteAddress(0)
	ppu.WriteSprite(10)
	ppu.WriteSprite('H')
	ppu.WriteSprite(0)
	ppu.WriteSprite(20)

	ppu.sprite.palettes[0][0] = color.RGBA{}
	assertScreen(t, ppu, 'H', 20, 10, ppu.sprite.palettes[0])
}

func TestPPUStatus(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))

	ppu := New(m)

	if !bits.IsFlag(ppu.Status(), 7) {
		t.Error("on vblank")
	}

	ppu.rendering = true
	if bits.IsFlag(ppu.Status(), 7) {
		t.Error("on rendering")
	}
}

func TestMirror(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	m.Write(0x2000, 'H')
	ppu := New(m)
	ppu.SetControl2(0x8)

	assertScreen(t, ppu, 'H', 0, 0, ppu.bg.palettes[0])
	assertScreen(t, ppu, ' ', 8, 0, ppu.bg.palettes[0])
}

func TestScroll(t *testing.T) {
	m := memory.New()
	m.Load(0x0, load("../example/hello/hello.nes"))
	m.Write(0x2001, 'H')
	ppu := New(m)
	ppu.SetControl2(0x8)

	ppu.SetScroll(1)
	ppu.SetScroll(0)
	assertScreen(t, ppu, 'H', 7, 0, ppu.bg.palettes[0])
}
