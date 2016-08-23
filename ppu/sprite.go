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
	copy(self.palettes[:], palettes)
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
			if sp.Y < HEIGHT {
				patterns[self.index][sp.Pattern].Put(
					img,
					image.Point{
						int(sp.X),
						int(sp.Y),
					},
					self.palettes[sp.Palette],
					pattern.Option{
						FlipH:         sp.FlipHorizon,
						FlipV:         sp.FlipVertical,
						BackdropColor: nil,
					})
			}
		}
	}
}

func (self *sprite) copyDMA(data []byte) {
	self.memory.Copy(data)
}
