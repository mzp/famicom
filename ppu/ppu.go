package ppu

import (
	"image"
	"image/color"

	memlib "github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/pattern"
)

type PPU struct {
	memory   *memlib.Memory
	patterns [2][]pattern.Pattern
}

func New(m *memlib.Memory) *PPU {
	t := PPU{memory: m}
	t.patterns[0] = pattern.ReadAllFromBytes(m.ReadRange(0x0, 0x1000))
	t.patterns[1] = pattern.ReadAllFromBytes(m.ReadRange(0x1000, 0x1000))
	return &t
}

func (ppu *PPU) screen() []byte {
	return ppu.memory.ReadRange(0x2000, 0x3C0)
}

func (ppu *PPU) Render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 240))

	// dummy color pallets
	pallets := []color.Color{
		color.RGBA{0, 0, 0, 0xFF},
		color.RGBA{0xA0, 0xA0, 0xA0, 0xFF},
		color.RGBA{0, 0xFF, 0, 0xFF},
		color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
	}

	for n, v := range ppu.screen() {
		x := n % 32
		y := n / 32

		if v != 0 {
			pattern.PutImage(img, x*8, y*8, ppu.patterns[0][v], pallets)
		}
	}

	return img
}
