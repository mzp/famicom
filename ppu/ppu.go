package ppu

import (
	"image"
	"image/color"

	"github.com/mzp/famicom/bits"
	memlib "github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/palette"
	"github.com/mzp/famicom/pattern"
)

type PPU struct {
	patterns [2][]pattern.Pattern

	sprite *sprite
	bg     bg
	vram   *vram

	rendering    bool
	interruptNMI bool
}

var black = color.RGBA{0, 0, 0, 0xFF}

func New(m *memlib.Memory) *PPU {
	t := PPU{}
	t.vram = makeVRAM(m)
	t.sprite = makeSprite()

	t.refresh()
	return &t
}

func (t *PPU) refresh() {
	m := t.vram.memory

	t.patterns[0] = pattern.ReadAllFromBytes(m.ReadRange(0x0, 0x1000))
	t.patterns[1] = pattern.ReadAllFromBytes(m.ReadRange(0x1000, 0x1000))

	t.bg.setPalettes([][]color.Color{
		palette.Read(m.ReadRange(0x3F00, 4)),
		palette.Read(m.ReadRange(0x3F04, 4)),
		palette.Read(m.ReadRange(0x3F08, 4)),
		palette.Read(m.ReadRange(0x3F0C, 4)),
	})

	t.sprite.setPalettes([][]color.Color{
		palette.Read(m.ReadRange(0x3F10, 4)),
		palette.Read(m.ReadRange(0x3F14, 4)),
		palette.Read(m.ReadRange(0x3F18, 4)),
		palette.Read(m.ReadRange(0x3F1C, 4)),
	})
}

func (ppu *PPU) startRender() {
	ppu.refresh()
	ppu.rendering = true
	ppu.sprite.hit = false
}

func (ppu *PPU) endRender() {
	ppu.rendering = false
}

func (ppu *PPU) Render() image.Image {
	ppu.startRender()
	defer ppu.endRender()

	img := ppu.bg.render(ppu.vram.memory, ppu.patterns)
	ppu.sprite.render(img, ppu.patterns)
	return img
}

func (ppu *PPU) Status() byte {
	ppu.bg.scroll.Reset()

	var n byte
	if ppu.sprite.hit {
		n = 0x40
	}

	if ppu.rendering {
		return 0x00 | n
	} else {
		return 0x80 | n
	}
}

func (ppu *PPU) SetControl1(flag byte) {
	ppu.bg.setNameTableIndex(int(flag & 0x3))

	ppu.vram.setOffset(bits.IsFlag(flag, 2))

	ppu.sprite.setIndex(bits.IsFlag(flag, 3))
	ppu.bg.setIndex(bits.IsFlag(flag, 4))

	ppu.interruptNMI = bits.IsFlag(flag, 7)
}

func (ppu *PPU) SetControl2(flag byte) {
	ppu.bg.enable = bits.IsFlag(flag, 3)
	ppu.sprite.enable = bits.IsFlag(flag, 4)
}

func (ppu *PPU) SetSpriteAddress(address uint8) {
	ppu.sprite.setAddress(address)
}

func (ppu *PPU) WriteSprite(value byte) {
	ppu.sprite.write(value)
}

func (ppu *PPU) SetAddress(data uint8) {
	ppu.vram.setAddress(data)
}

func (ppu *PPU) WriteVRAM(data uint8) {
	ppu.vram.write(data)
}

func (ppu *PPU) ReadVRAM() uint8 {
	return ppu.vram.read()
}

func (ppu *PPU) CopySpriteDMA(data []byte) {
	ppu.sprite.copyDMA(data)
}

func (ppu *PPU) IsInterrupNMI() bool {
	return ppu.interruptNMI
}

func (ppu *PPU) SetScroll(value byte) {
	ppu.bg.setScroll(value)
}

func (ppu *PPU) SetVerticalMirror(mirror bool) {
	ppu.bg.verticalMirror = mirror
}
