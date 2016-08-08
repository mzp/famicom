package ioregister

import (
	"testing"

	"github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/ppu"
)

func TestConnectPPU(t *testing.T) {
	cpuMemory := memory.New()
	ppuMemory := memory.New()
	p := ppu.New(ppuMemory)

	ConnectPPU(cpuMemory, p)

	cpuMemory.Write(0x2006, 0x20)
	cpuMemory.Write(0x2006, 0x00)
	cpuMemory.Write(0x2007, 0x42)

	if ppuMemory.Read(0x2000) != 0x42 {
		t.Errorf("cannot write value: %x", ppuMemory.Read(0x2000))
	}
}
