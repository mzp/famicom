package ppu

import (
	"testing"

	"github.com/mzp/famicom/memory"
)

func TestSetAddress(t *testing.T) {
	m := memory.New()
	vram := makeVRAM(m)
	vram.setAddress(0x20)
	vram.setAddress(0x0)
	if vram.address.Value() != 0x2000 {
		t.Errorf("expect 0x2000 but %x", vram.address)
	}

	vram.setAddress(0xca)
	vram.setAddress(0xfe)
	if vram.address.Value() != 0xcafe {
		t.Errorf("expect 0xcafe but %x", vram.address)
	}
}
