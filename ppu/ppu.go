package ppu

import (
	"image"
	"image/color"

	"github.com/mzp/famicom/bits"
	memlib "github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/palette"
	"github.com/mzp/famicom/pattern"
	"github.com/mzp/famicom/sprite"
)

type PPU struct {
	memory          *memlib.Memory
	spriteMemory    *sprite.SpriteMemory
	patterns        [2][]pattern.Pattern
	spriteIndex     int
	backgroundIndex int
	bgPalettes      [4][]color.Color
	spritePalettes  [4][]color.Color
	nameTable       uint16
	sprite          bool
	background      bool
	vramAddress     uint16
	vramHigh        bool
	vramOffset      uint16
	rendering       bool
}

var black = color.RGBA{0, 0, 0, 0xFF}

func New(m *memlib.Memory) *PPU {
	t := PPU{}
	t.memory = m
	t.spriteMemory = sprite.New()

	t.nameTable = 0x2000
	t.vramHigh = true
	t.vramOffset = 1

	t.refresh()
	return &t
}

func (t *PPU) refresh() {
	m := t.memory
	t.patterns[0] = pattern.ReadAllFromBytes(m.ReadRange(0x0, 0x1000))
	t.patterns[1] = pattern.ReadAllFromBytes(m.ReadRange(0x1000, 0x1000))

	t.bgPalettes[0] = palette.Read(m.ReadRange(0x3F00, 4))
	t.bgPalettes[1] = palette.Read(m.ReadRange(0x3F04, 4))
	t.bgPalettes[2] = palette.Read(m.ReadRange(0x3F08, 4))
	t.bgPalettes[3] = palette.Read(m.ReadRange(0x3F0C, 4))

	t.spritePalettes[0] = palette.Read(m.ReadRange(0x3F10, 4))
	t.spritePalettes[1] = palette.Read(m.ReadRange(0x3F14, 4))
	t.spritePalettes[2] = palette.Read(m.ReadRange(0x3F18, 4))
	t.spritePalettes[3] = palette.Read(m.ReadRange(0x3F1C, 4))
}

func (ppu *PPU) SetControl1(flag byte) {
	ppu.nameTable = 0x2000 + 0x400*uint16(flag&0x3)

	if bits.IsFlag(flag, 2) {
		ppu.vramOffset = 32
	} else {
		ppu.vramOffset = 1
	}

	if bits.IsFlag(flag, 3) {
		ppu.spriteIndex = 1
	} else {
		ppu.spriteIndex = 0
	}

	if bits.IsFlag(flag, 4) {
		ppu.backgroundIndex = 1
	} else {
		ppu.backgroundIndex = 0
	}
}

func (ppu *PPU) SetControl2(flag byte) {
	ppu.background = bits.IsFlag(flag, 3)
	ppu.sprite = bits.IsFlag(flag, 4)
}

func (ppu *PPU) SetSpriteAddress(address uint8) {
	ppu.spriteMemory.SetAddress(address)
}

func (ppu *PPU) WriteSprite(value byte) {
	ppu.spriteMemory.Write(value)
}

func (ppu *PPU) SetAddress(data uint8) {
	if ppu.vramHigh {
		ppu.vramAddress = uint16(data) << 8
		ppu.vramHigh = false
	} else {
		ppu.vramAddress |= uint16(data)
	}
}

func (ppu *PPU) WriteVRAM(data uint8) {
	ppu.memory.Write(ppu.vramAddress, data)
	ppu.vramAddress += ppu.vramOffset
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

func (ppu *PPU) startRender() {
	ppu.refresh()
	ppu.rendering = true
}

func (ppu *PPU) endRender() {
	ppu.rendering = false
}

func (ppu *PPU) Render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 240))
	ppu.startRender()
	defer ppu.endRender()

	nameTable, attributeTable := ppu.screen()
	for n, v := range nameTable {
		x := n % 32
		y := n / 32

		if ppu.background {
			paletteIndex := getAttribute(attributeTable, x, y)
			pattern.PutImage(img,
				x*8, y*8,
				ppu.patterns[ppu.backgroundIndex][v],
				ppu.bgPalettes[paletteIndex])
		}
	}

	if ppu.sprite {
		for _, sp := range ppu.spriteMemory.Get() {
			pattern.PutImage(img,
				int(sp.X), int(sp.Y),
				ppu.patterns[ppu.spriteIndex][sp.Pattern],
				ppu.spritePalettes[sp.Palette])
		}
	}

	return img
}

func (ppu *PPU) Status() byte {
	if ppu.rendering {
		return 0
	} else {
		return 0x80
	}
}
