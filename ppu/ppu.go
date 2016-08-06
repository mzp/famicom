package ppu

import (
	"image"
	"image/color"

	memlib "github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/palette"
	"github.com/mzp/famicom/pattern"
)

type PPU struct {
	memory     *memlib.Memory
	patterns   [2][]pattern.Pattern
	bgPalettes [4][]color.Color
	nameTable  uint16
}

func New(m *memlib.Memory) *PPU {
	t := PPU{memory: m}
	t.patterns[0] = pattern.ReadAllFromBytes(m.ReadRange(0x0, 0x1000))
	t.patterns[1] = pattern.ReadAllFromBytes(m.ReadRange(0x1000, 0x1000))

	t.bgPalettes[0] = palette.Read(m.ReadRange(0x3F00, 4))
	t.bgPalettes[1] = palette.Read(m.ReadRange(0x3F04, 4))
	t.bgPalettes[2] = palette.Read(m.ReadRange(0x3F08, 4))
	t.bgPalettes[3] = palette.Read(m.ReadRange(0x3F0C, 4))

	t.nameTable = 0x2000
	return &t
}

func (ppu *PPU) SetControl1(flag byte) {
	ppu.nameTable = 0x2000 + 0x400*uint16(flag&0x3)
}

func (ppu *PPU) SetControl2(flag byte) {
}

func (ppu *PPU) screen() ([]byte, []byte) {
	const (
		NameTableSize      = 0x3c0
		AttributeTableSize = 0x40
	)

	nameTable := ppu.memory.ReadRange(ppu.nameTable, NameTableSize)
	attributeTable := ppu.memory.ReadRange(ppu.nameTable+NameTableSize, AttributeTableSize)

	return nameTable, attributeTable
}

func getAttribute(attributeTable []byte, x, y int) byte {
	attribute := attributeTable[x/4+y/4*8]

	x_, y_ := x%4, y%4
	index := (x_ / 2) + (y_/2)*2

	return (attribute >> uint(index*2)) & 0x3
}

func (ppu *PPU) Render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 240))

	nameTable, attributeTable := ppu.screen()
	for n, v := range nameTable {
		x := n % 32
		y := n / 32

		paletteIndex := getAttribute(attributeTable, x, y)
		pattern.PutImage(img, x*8, y*8, ppu.patterns[0][v], ppu.bgPalettes[paletteIndex])
	}

	return img
}
