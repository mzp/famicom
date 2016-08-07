package sprite

import (
	"github.com/mzp/famicom/bits"
	"github.com/mzp/famicom/memory"
)

type SpriteMemory struct {
	memory  memory.Memory
	address uint8
}

type Sprite struct {
	Palette         uint
	Pattern         uint
	FlipVertical    bool
	FlipHorizon     bool
	UnderBackground bool
	X, Y            uint
}

func New() *SpriteMemory {
	t := SpriteMemory{}
	return &t
}

func (t *SpriteMemory) SetAddress(address uint8) {
	t.address = address
}

func (t *SpriteMemory) Write(value byte) {
	t.memory.Write(uint16(t.address), value)
	t.address += 1
}

func parse(data []byte) Sprite {
	return Sprite{
		Y:               uint(data[0]),
		Pattern:         uint(data[1]),
		FlipVertical:    bits.IsFlag(data[2], 7),
		FlipHorizon:     bits.IsFlag(data[2], 6),
		UnderBackground: bits.IsFlag(data[2], 5),
		Palette:         uint(data[2] & 0x3),
		X:               uint(data[3]),
	}
}

func (t *SpriteMemory) Get() []Sprite {
	var xs [64]Sprite

	for i := uint16(0); i < 64; i++ {
		xs[i] = parse(t.memory.ReadRange(i*4, 4))
	}

	return xs[:]
}
