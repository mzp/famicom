package ioregister

import (
	"fmt"

	"github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/pad"
	"github.com/mzp/famicom/ppu"
)

func ConnectPPU(memory *memory.Memory, ppu *ppu.PPU) {
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
	memory.WriteTrap(0x2008, func(value byte) {
		fmt.Println(value)
	})

	memory.WriteTrap(0x4014, func(value byte) {
		data := memory.ReadRange(uint16(value) << 8, 0xFF)
		ppu.CopySpriteDMA(data)
	})
}

func ConnectPad(memory *memory.Memory, pad1, pad2 *pad.Pad) {
	memory.ReadTrap(0x4016, func() byte {
		return pad1.Read()
	})
	memory.WriteTrap(0x4016, func(value byte) {
		pad1.Write(value)
	})

	memory.ReadTrap(0x4017, func() byte {
		return pad2.Read()
	})
	memory.WriteTrap(0x4017, func(value byte) {
		pad2.Write(value)
	})
}
