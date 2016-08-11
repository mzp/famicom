package cpu

import "testing"

func TestReset(t *testing.T) {
	cpu, memory := create()
	memory.Write(0xFFFC, 0xfe)
	memory.Write(0xFFFD, 0xca)
	cpu.InterrruptReset()

	cpu.Step()

	if cpu.pc != 0xcafe {
		t.Error("cannot jump Reset handler")
	}

	if !cpu.status.irq {
		t.Error("cannot disable irq status")
	}
}

func TestNMI(t *testing.T) {
	cpu, memory := create()
	memory.Write(0xFFFA, 0xfe)
	memory.Write(0xFFFB, 0xca)
	cpu.status.negative = true
	cpu.status.brk = true

	cpu.InterrruptNMI()

	cpu.Step()

	if cpu.pc != 0xcafe {
		t.Error("cannot jump Reset handler")
	}

	if !cpu.status.irq {
		t.Error("cannot disable irq status")
	}

	if cpu.status.brk {
		t.Error("cannot reset break status")
	}

	if memory.Read16(0x01FE) != 0 {
		t.Errorf("cannot push current pc. %x", memory.Read16(0x01FE))
	}

	if memory.Read(0x01FD) != 0x90 {
		t.Errorf("cannot status: %x", memory.Read(0x01FD))
	}

	if cpu.s != 0xFC {
		t.Errorf("cannot push down stack: %x", cpu.s)
	}
}

func TestIrq(t *testing.T) {
	cpu, memory := create()
	memory.Write(0xFFFE, 0xfe)
	memory.Write(0xFFFF, 0xca)
	cpu.status.negative = true
	cpu.status.brk = true

	cpu.InterrruptIrq()

	cpu.Step()

	if cpu.pc != 0xcafe {
		t.Error("cannot jump Reset handler")
	}

	if !cpu.status.irq {
		t.Error("cannot disable irq status")
	}

	if cpu.status.brk {
		t.Error("cannot reset break status")
	}

	if memory.Read16(0x01FE) != 0 {
		t.Errorf("cannot push current pc. %x", memory.Read16(0x01FE))
	}

	if memory.Read(0x01FD) != 0x90 {
		t.Errorf("cannot status: %x", memory.Read(0x01FD))
	}

	if cpu.s != 0xFC {
		t.Errorf("cannot push down stack: %x", cpu.s)
	}
}

func TestBreak(t *testing.T) {
	cpu, memory := create()
	memory.Write(0xFFFE, 0xfe)
	memory.Write(0xFFFF, 0xca)
	cpu.status.negative = true
	cpu.status.brk = true

	cpu.InterrruptBreak()

	cpu.Step()

	if cpu.pc != 0xcafe {
		t.Error("cannot jump Reset handler")
	}

	if !cpu.status.irq {
		t.Error("cannot disable irq status")
	}

	if !cpu.status.brk {
		t.Error("cannot set break status")
	}

	if memory.Read16(0x01FE) != 0 {
		t.Errorf("cannot push current pc. %x", memory.Read16(0x01FE))
	}

	if memory.Read(0x01FD) != 0x90 {
		t.Errorf("cannot status: %x", memory.Read(0x01FD))
	}

	if cpu.s != 0xFC {
		t.Errorf("cannot push down stack: %x", cpu.s)
	}
}

func TestIrqDisable(t *testing.T) {
	cpu, memory := create()
	memory.Write(0xFFFE, 0xfe)
	memory.Write(0xFFFF, 0xca)
	cpu.status.irq = true

	cpu.InterrruptIrq()

	cpu.Step()

	if cpu.pc == 0xcafe {
		t.Error("cannot disable irq")
	}
}

func TestBreakDisable(t *testing.T) {
	cpu, memory := create()
	memory.Write(0xFFFE, 0xfe)
	memory.Write(0xFFFF, 0xca)
	cpu.status.irq = true

	cpu.InterrruptBreak()

	cpu.Step()

	if cpu.pc == 0xcafe {
		t.Error("cannot disable irq")
	}
}
