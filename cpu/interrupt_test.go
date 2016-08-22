package cpu

import "testing"

func TestReset(t *testing.T) {
	core, memory := createCore()
	memory.Write(0xFFFC, 0xfe)
	memory.Write(0xFFFD, 0xca)
	handleReset(core)

	if core.pc != 0xcafe {
		t.Error("cannot jump Reset handler")
	}

	if !core.status.irq {
		t.Error("cannot disable irq status")
	}
}

func TestNMI(t *testing.T) {
	core, memory := createCore()
	memory.Write(0xFFFA, 0xfe)
	memory.Write(0xFFFB, 0xca)
	core.status.negative = true
	core.status.brk = true

	handleNMI(core)

	if core.pc != 0xcafe {
		t.Error("cannot jump Reset handler")
	}

	if !core.status.irq {
		t.Error("cannot disable irq status")
	}

	if core.status.brk {
		t.Error("cannot reset break status")
	}

	if memory.Read16(0x01FE) != 0 {
		t.Errorf("cannot push current pc. %x", memory.Read16(0x01FE))
	}

	if memory.Read(0x01FD) != 0x90 {
		t.Errorf("cannot status: %x", memory.Read(0x01FD))
	}

	if core.s != 0xFC {
		t.Errorf("cannot push down stack: %x", core.s)
	}
}

func TestIrq(t *testing.T) {
	core, memory := createCore()
	memory.Write(0xFFFE, 0xfe)
	memory.Write(0xFFFF, 0xca)
	core.status.negative = true
	core.status.brk = true

	handleInteruption(core, false)

	if core.pc != 0xcafe {
		t.Error("cannot jump Reset handler")
	}

	if !core.status.irq {
		t.Error("cannot disable irq status")
	}

	if core.status.brk {
		t.Error("cannot reset break status")
	}

	if memory.Read16(0x01FE) != 0 {
		t.Errorf("cannot push current pc. %x", memory.Read16(0x01FE))
	}

	if memory.Read(0x01FD) != 0x90 {
		t.Errorf("cannot status: %x", memory.Read(0x01FD))
	}

	if core.s != 0xFC {
		t.Errorf("cannot push down stack: %x", core.s)
	}
}

func TestBreak(t *testing.T) {
	core, memory := createCore()
	memory.Write(0xFFFE, 0xfe)
	memory.Write(0xFFFF, 0xca)
	core.status.negative = true
	core.status.brk = true

	handleInteruption(core, true)

	if core.pc != 0xcafe {
		t.Error("cannot jump Reset handler")
	}

	if !core.status.irq {
		t.Error("cannot disable irq status")
	}

	if !core.status.brk {
		t.Error("cannot set break status")
	}

	if memory.Read16(0x01FE) != 0 {
		t.Errorf("cannot push current pc. %x", memory.Read16(0x01FE))
	}

	if memory.Read(0x01FD) != 0x90 {
		t.Errorf("cannot status: %x", memory.Read(0x01FD))
	}

	if core.s != 0xFC {
		t.Errorf("cannot push down stack: %x", core.s)
	}
}
