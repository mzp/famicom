package ppu

import (
	"image"
	"image/color"

	"github.com/mzp/famicom/debug"
	"github.com/mzp/famicom/pattern"
	s "github.com/mzp/famicom/sprite"
)

type sprite struct {
	memory   *s.SpriteMemory
	enable   bool
	index    int
	palettes [4][]color.Color
	hit      bool
}

func makeSprite() *sprite {
	t := sprite{}
	t.memory = s.New()
	return &t
}

func (self *sprite) setPalettes(palettes [][]color.Color) {
	self.palettes[0] = palettes[0]
	self.palettes[1] = palettes[1]
	self.palettes[2] = palettes[2]
	self.palettes[3] = palettes[3]
}

func (self *sprite) setIndex(upper bool) {
	if upper {
		self.index = 1
	} else {
		self.index = 0
	}
}

func (self *sprite) setAddress(address uint8) {
	self.memory.SetAddress(address)
}

func (self *sprite) write(value byte) {
	self.memory.Write(value)
}

func (self *sprite) render(img *image.RGBA, patterns [2][]pattern.Pattern) {
	debug.DumpSprite(self.memory.Get())

	// TODO: mock implement
	self.hit = true

	if self.enable {
		for _, sp := range self.memory.Get() {
			pattern.PutImageWithFlip(img,
				int(sp.X),
				int(sp.Y),
				patterns[self.index][sp.Pattern],
				self.palettes[sp.Palette],
				sp.FlipHorizon,
				sp.FlipVertical,
				true)
		}
	}
}

func (self *sprite) copyDMA(data []byte) {
	self.memory.Copy(data)
}
