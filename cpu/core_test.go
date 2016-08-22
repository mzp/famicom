package cpu

import (
	"github.com/mzp/famicom/decoder"
	"github.com/mzp/famicom/memory"
	"testing"
)

func createCore() (*core, *memory.Memory) {
	m := memory.New()
	core := makeCore(m, 0)
	return core, m
}

func TestCoreRead(t *testing.T) {
	core, m := createCore()
	m.Write(0x01, 0x00)
	m.Write(0x02, 0x20)

	m.Write(0x80, 0x42)
	m.Write(0x81, 0x43)
	m.Write(0x82, 0x44)

	m.Write(0x2000, 0x80)
	m.Write(0x2001, 0x81)
	m.Write(0x2002, 0x82)

	core.x = 1
	core.y = 2

	tests := []struct {
		mode   decoder.AddressingMode
		value  int
		expect uint8
	}{
		{decoder.Immediate, 1, 1},
		{decoder.ZeroPage, 0x80, 0x42},
		{decoder.ZeroPageX, 0x80, 0x43},
		{decoder.ZeroPageY, 0x80, 0x44},
		{decoder.Absolute, 0x2000, 0x80},
		{decoder.AbsoluteX, 0x2000, 0x81},
		{decoder.AbsoluteY, 0x2000, 0x82},
		{decoder.IndirectX, 0x00, 0x80},
		{decoder.IndirectY, 0x01, 0x82},
	}

	for _, test := range tests {
		inst := decoder.Instruction{
			AddressingMode: test.mode,
			Value:          test.value,
		}
		value := core.read(inst)

		if value != test.expect {
			t.Errorf("%s %x != %x", inst.String(), value, test.expect)
		}
	}
}

func TestNZ(t *testing.T) {
	core, _ := createCore()

	tests := []struct {
		value  uint8
		status status
	}{
		{0, status{zero: true}},
		{1, status{}},
		{0x80, status{negative: true}},
	}

	for _, test := range tests {
		core.nz(test.value)
		if core.status != test.status {
			t.Errorf("%s != %s", core.status, test.status)
		}
	}
}

func TestLoad(t *testing.T) {
	core, _ := createCore()
	core.status.overflow = true
	core.status.carry = true
	core.status.brk = true

	core.load(&core.a, 42)

	if core.a != 42 {
		t.Errorf("cpu.a = %x, but must be 42", core.a)
	}

	if (core.status != status{overflow: true, carry: true, brk: true}) {
		t.Error("clear unexpected flag")
	}
}

func TestStore(t *testing.T) {
	core, m := createCore()

	core.status.zero = true
	core.status.overflow = true
	core.status.carry = true
	core.status.brk = true

	core.store(42, 1)
	if m.Read(42) != 1 {
		t.Errorf("Store 1, but %x", m.Read(42))
	}

	if (core.status != status{overflow: true, carry: true, brk: true, zero: true}) {
		t.Error("clear unexpected flag")
	}
}
