package ioregister

import (
	"github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/ppu"
)

func Connect(memory *memory.Memory, ppu *ppu.PPU) {
	memory.WriteTrap(0x2000, func(value byte) {
		ppu.SetControl1(value)
	})
	memory.WriteTrap(0x2001, func(value byte) {
		ppu.SetControl2(value)
	})
	memory.ReadTrap(0x2002, func() byte {
		return ppu.Status()
	})
	memory.WriteTrap(0x2003, func(value byte) {
		ppu.SetSpriteAddress(value)
	})
	memory.WriteTrap(0x2004, func(value byte) {
		ppu.WriteSprite(value)
	})

	memory.WriteTrap(0x2006, func(value byte) {
		ppu.SetAddress(value)
	})
	memory.WriteTrap(0x2007, func(value byte) {
		ppu.WriteVRAM(value)
	})
}
