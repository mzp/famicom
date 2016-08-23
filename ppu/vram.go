package ppu

import (
	"github.com/mzp/famicom/memory"
)

type vram struct {
	memory  *memory.Memory
	address doubleByte
	buffer  uint8
	offset  uint16
}

func makeVRAM(m *memory.Memory) *vram {
	t := vram{memory: m, offset: 1}

	m.SetMirror(0x3F00, 0x3F10, 1)
	m.SetMirror(0x3F04, 0x3F14, 1)
	m.SetMirror(0x3F08, 0x3F18, 1)
	m.SetMirror(0x3F0C, 0x3F1C, 1)

	return &t
}

func (self *vram) setAddress(data uint8) {
	self.address.Write(data)
}

func (self *vram) write(data uint8) {
	address := self.address.Value()
	defer self.address.Set(address + self.offset)

	self.memory.Write(address, data)
}

func (self *vram) read() uint8 {
	address := self.address.Value()
	defer self.address.Set(address + self.offset)

	value := self.memory.Read(address)

	if address < 0x3F00 {
		t := self.buffer
		self.buffer = value
		return t
	} else {
		return value
	}
}

func (self *vram) setOffset(largeOffset bool) {
	if largeOffset {
		self.offset = 32
	} else {
		self.offset = 1
	}
}
