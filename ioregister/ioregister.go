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
	memory.WriteTrap(0x2006, func(value byte) {
		ppu.SetAddress(value)
	})
	memory.WriteTrap(0x2007, func(value byte) {
		ppu.WriteVRAM(value)
	})
}
