package ppu

import (
	"image"
	"image/color"

	"github.com/mzp/famicom/debug"
	"github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/pattern"
)

type bg struct {
	enable         bool
	verticalMirror bool
	palettes       [4][]color.Color
	origin         image.Point
	scroll         doubleByte
	index          int
}

func (self *bg) setPalettes(palettes [][]color.Color) {
	// FIXME: palette parser
	self.palettes[0] = palettes[0]
	self.palettes[1] = palettes[1]
	self.palettes[2] = palettes[2]
	self.palettes[3] = palettes[3]

	// Use universal background color at all palette 0
	for i := 0; i < 3; i++ {
		self.palettes[i+1][0] = self.palettes[0][0]
	}
}

func (self *bg) setNameTableIndex(n int) {
	self.origin.X = n % 2
	self.origin.Y = n / 2
}

func (self *bg) setIndex(upper bool) {
	if upper {
		self.index = 1
	} else {
		self.index = 0
	}
}

func (self *bg) setScroll(value byte) {
	self.scroll.Write(value)
}

func (self *bg) getTables(vram *memory.Memory) ([][]byte, [][]byte) {
	if self.verticalMirror {
		return [][]byte{
				vram.ReadRange(0x2000, 0x3c0),
				vram.ReadRange(0x2400, 0x3c0),
				vram.ReadRange(0x2000, 0x3c0),
				vram.ReadRange(0x2400, 0x3c0),
			}, [][]byte{
				vram.ReadRange(0x23c0, 0x40),
				vram.ReadRange(0x27c0, 0x40),
				vram.ReadRange(0x23c0, 0x40),
				vram.ReadRange(0x27c0, 0x40),
			}
	} else {
		return [][]byte{
				vram.ReadRange(0x2000, 0x3c0),
				vram.ReadRange(0x2000, 0x3c0),
				vram.ReadRange(0x2800, 0x3c0),
				vram.ReadRange(0x2800, 0x3c0),
			}, [][]byte{
				vram.ReadRange(0x23c0, 0x40),
				vram.ReadRange(0x23c0, 0x40),
				vram.ReadRange(0x2bc0, 0x40),
				vram.ReadRange(0x2bc0, 0x40),
			}
	}
}

func (self *bg) renderAll(vram *memory.Memory, patterns [2][]pattern.Pattern) *image.RGBA {
	background := image.NewRGBA(image.Rect(0, 0, WIDTH*2, HEIGHT*2))

	nameTables, attributeTables := self.getTables(vram)

	for i := 0; i < 4; i++ {
		nameTable := nameTables[i]
		attributeTable := attributeTables[i]

		var debugingAttributeTable [32 * 32]byte

		bx := (i % 2) * WIDTH
		by := (i / 2) * HEIGHT

		for n, v := range nameTable {
			x := n % 32
			y := n / 32

			paletteIndex := getAttribute(attributeTable, x, y)

			debugingAttributeTable[y*32+x] = paletteIndex
			pattern.PutImage(background,
				bx+x*8, by+y*8,
				patterns[self.index][v],
				self.palettes[paletteIndex])
		}
		debug.DumpBackground(i, nameTable, debugingAttributeTable[:])
	}

	return background
}

func (self *bg) render(vram *memory.Memory, patterns [2][]pattern.Pattern) *image.RGBA {
	if self.enable {
		background := self.renderAll(vram, patterns)

		return clip(background,
			self.origin.X*WIDTH+int(self.scroll.data[0]),
			self.origin.Y*HEIGHT+int(self.scroll.data[1]),
			WIDTH,
			HEIGHT)
	} else {
		return image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	}
}

func getAttribute(attributeTable []byte, x, y int) byte {
	attribute := attributeTable[x/4+y/4*8]

	x_, y_ := x%4, y%4
	index := (x_ / 2) + (y_/2)*2

	return (attribute >> uint(index*2)) & 0x3
}

func wrap(n, m int) int {
	if n >= 0 {
		return n % m
	} else {
		return m + n
	}
}

func clip(img *image.RGBA, x, y, width, height int) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, width, height))
	bounds := img.Bounds()
	originalWidth := bounds.Max.X - bounds.Min.X
	originalHeight := bounds.Max.Y - bounds.Min.Y

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			c := img.At(wrap(x+i, originalWidth), wrap(y+j, originalHeight))
			result.Set(i, j, c)
		}
	}
	return result
}
