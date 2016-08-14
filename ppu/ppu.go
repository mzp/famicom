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
	VerticalMirror   bool
	memory           *memlib.Memory
	spriteMemory     *sprite.SpriteMemory
	patterns         [2][]pattern.Pattern
	spriteIndex      int
	backgroundIndex  int
	interruptNMI     bool
	bgPalettes       [4][]color.Color
	spritePalettes   [4][]color.Color
	originX, originY int
	sprite           bool
	background       bool
	vramAddress      uint16
	vramHigh         bool
	vramOffset       uint16
	rendering        bool
}

var black = color.RGBA{0, 0, 0, 0xFF}

func New(m *memlib.Memory) *PPU {
	t := PPU{}
	t.memory = m
	t.spriteMemory = sprite.New()

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
	origin := int(flag & 0x3)
	ppu.originX = origin % 2
	ppu.originY = origin / 2

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

	ppu.interruptNMI = bits.IsFlag(flag, 7)
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

func (ppu *PPU) startRender() {
	ppu.refresh()
	ppu.rendering = true
}

func (ppu *PPU) endRender() {
	ppu.rendering = false
}

func (ppu *PPU) Render() image.Image {
	ppu.startRender()
	defer ppu.endRender()

	var img *image.RGBA
	if ppu.background {
		var background *image.RGBA

		if ppu.VerticalMirror {
			background = renderBackground(ppu,
				[4][]byte{
					ppu.memory.ReadRange(0x2000, 0x3c0),
					ppu.memory.ReadRange(0x2400, 0x3c0),
					ppu.memory.ReadRange(0x2000, 0x3c0),
					ppu.memory.ReadRange(0x2400, 0x3c0),
				},
				[4][]byte{
					ppu.memory.ReadRange(0x23c0, 0x40),
					ppu.memory.ReadRange(0x27c0, 0x40),
					ppu.memory.ReadRange(0x23c0, 0x40),
					ppu.memory.ReadRange(0x27c0, 0x40),
				})
		} else {
			background = renderBackground(ppu,
				[4][]byte{
					ppu.memory.ReadRange(0x2000, 0x3c0),
					ppu.memory.ReadRange(0x2000, 0x3c0),
					ppu.memory.ReadRange(0x2800, 0x3c0),
					ppu.memory.ReadRange(0x2800, 0x3c0),
				},
				[4][]byte{
					ppu.memory.ReadRange(0x23c0, 0x40),
					ppu.memory.ReadRange(0x23c0, 0x40),
					ppu.memory.ReadRange(0x2bc0, 0x40),
					ppu.memory.ReadRange(0x2bc0, 0x40),
				})
		}
		img = background.SubImage(image.Rect(
			ppu.originX*WIDTH,
			ppu.originY*HEIGHT,
			(ppu.originX+1)*WIDTH,
			(ppu.originY+1)*HEIGHT)).(*image.RGBA)
	} else {
		img = image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	}

	if ppu.sprite {
		for _, sp := range ppu.spriteMemory.Get() {
			pattern.PutImage(img,
				img.Rect.Min.X+int(sp.X),
				img.Rect.Min.Y+int(sp.Y),
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

func (ppu *PPU) CopySpriteDMA(data []byte) {
	ppu.spriteMemory.Copy(data)
}

func (ppu *PPU) IsInterrupNMI() bool {
	return ppu.interruptNMI
}
