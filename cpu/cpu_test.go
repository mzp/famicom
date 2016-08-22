package cpu

import (
	"github.com/mzp/famicom/decoder"
	"github.com/mzp/famicom/memory"
	"testing"
)

func create() (*CPU, *memory.Memory) {
	m := memory.New()
	cpu := New(m, 0)
	return cpu, m
}

func TestTransfer(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.None,
			Value:          0,
		}
	}

	cpu, _ := create()

	cpu.core.a = 1
	cpu.Execute(inst(decoder.TAX))
	if cpu.core.x != 1 {
		t.Errorf("Transfer 1, but %x", cpu.core.x)
	}

	cpu.Execute(inst(decoder.TAY))
	if cpu.core.y != 1 {
		t.Errorf("Transfer 1, but %x", cpu.core.y)
	}

	cpu.core.s = 2
	cpu.Execute(inst(decoder.TSX))
	if cpu.core.x != 2 {
		t.Errorf("Transfer 2, but %x", cpu.core.x)
	}

	cpu.core.x = 3
	cpu.Execute(inst(decoder.TXA))
	if cpu.core.a != 3 {
		t.Errorf("Transfer 3, but %x", cpu.core.a)
	}

	cpu.Execute(inst(decoder.TXS))
	if cpu.core.s != 3 {
		t.Errorf("Transfer 3, but %x", cpu.core.s)
	}

	cpu.core.y = 4
	cpu.Execute(inst(decoder.TYA))
	if cpu.core.a != 4 {
		t.Errorf("Transfer 4, but %x", cpu.core.a)
	}
}

func TestADC(t *testing.T) {
	inst := decoder.Instruction{
		Op:             decoder.ADC,
		AddressingMode: decoder.Absolute,
		Value:          0x2000,
	}

	cpu, m := create()

	tests := []struct {
		a      uint8
		memory uint8
		carry  bool
		expect uint8
		status status
	}{
		{0x30, 0x20, false, 0x50, status{}},
		{0x0, 0x0, false, 0x0, status{zero: true}},
		{0x0, 0x0, true, 0x1, status{}},
		{0x7F, 0x1, false, 0x80, status{negative: true, overflow: true}},
		{0x80, 0x5, false, 0x85, status{negative: true}},
		{0x80, 0x80, false, 0x00, status{overflow: true, carry: true, zero: true}},
	}

	for _, test := range tests {
		cpu.core.a = test.a
		cpu.core.status.carry = test.carry
		m.Write(0x2000, test.memory)

		cpu.Execute(inst)

		if cpu.core.a != test.expect {
			t.Errorf("Expect %x, but %x", test.expect, cpu.core.a)
		}

		if cpu.core.status != test.status {
			t.Errorf("Expect\n%s, but\n%s", test.status.String(), cpu.core.status.String())
		}
	}
}

func TestADCFlag(t *testing.T) {
	inst := decoder.Instruction{
		Op:             decoder.ADC,
		AddressingMode: decoder.Immediate,
		Value:          0x20,
	}

	cpu, _ := create()
	cpu.core.a = 0x30
	cpu.core.status.brk = true

	cpu.Execute(inst)

	if !cpu.core.status.brk {
		t.Error("clear unexpected flag")
	}
}

func TestSBC(t *testing.T) {
	inst := decoder.Instruction{
		Op:             decoder.SBC,
		AddressingMode: decoder.Absolute,
		Value:          0x2000,
	}

	cpu, m := create()

	tests := []struct {
		a      uint8
		memory uint8
		carry  bool
		expect uint8
		status status
	}{
		{0x30, 0x20, true, 0x10, status{carry: true}},
		{0x0, 0x0, true, 0x0, status{zero: true, carry: true}},
		{0x2, 0x0, false, 0x1, status{carry: true}},
		{0x0, 0x1, true, 0xFF, status{negative: true, overflow: true}},
		{0x1, 0x3, true, 0xFE, status{carry: false, overflow: true, negative: true}},
		{0x80, 0x1, true, 0x7F, status{overflow: true, carry: true}},
		{0x85, 0x1, true, 0x84, status{negative: true, carry: true}},
	}

	for i, test := range tests {
		cpu.core.a = test.a
		cpu.core.status.carry = test.carry
		m.Write(0x2000, test.memory)

		cpu.Execute(inst)

		if cpu.core.a != test.expect {
			t.Errorf("Expect %x, but %x", test.expect, cpu.core.a)
		}

		if cpu.core.status != test.status {
			t.Errorf("%d. Expect\n%s, but\n%s", i, test.status.String(), cpu.core.status.String())
		}
	}
}

func TestSBCFlag(t *testing.T) {
	inst := decoder.Instruction{
		Op:             decoder.SBC,
		AddressingMode: decoder.Immediate,
		Value:          0x20,
	}

	cpu, _ := create()
	cpu.core.a = 0x30
	cpu.core.status.brk = true

	cpu.Execute(inst)

	if !cpu.core.status.brk {
		t.Error("clear unexpected flag")
	}
}

func TestBitOp(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.Absolute,
			Value:          0x2000,
		}
	}

	cpu, m := create()
	m.Write(0x2000, 0xde)
	cpu.core.a = 0xad
	cpu.Execute(inst(decoder.AND))
	if cpu.core.a != 0x8c {
		t.Errorf("0xde & 0xad = %x", cpu.core.a)
	}

	cpu.core.a = 0xad
	cpu.Execute(inst(decoder.EOR))
	if cpu.core.a != 0x73 {
		t.Errorf("0xde ^ 0xad = %x", cpu.core.a)
	}

	cpu.core.a = 0xad
	cpu.Execute(inst(decoder.ORA))
	if cpu.core.a != 0xff {
		t.Errorf("0xde | 0xad = %x", cpu.core.a)
	}
}

func TestIncDec(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.Absolute,
			Value:          0x2000,
		}
	}

	cpu, m := create()
	m.Write(0x2000, 0x80)

	cpu.Execute(inst(decoder.INC))
	if m.Read(0x2000) != 0x81 {
		t.Error()
	}

	cpu.Execute(inst(decoder.DEC))
	if m.Read(0x2000) != 0x80 {
		t.Error()
	}

	if !cpu.core.status.negative {
		t.Error()
	}
}

func TestIncDecReg(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.None,
			Value:          0,
		}
	}

	cpu, _ := create()

	cpu.Execute(inst(decoder.INX))
	if cpu.core.x != 1 {
		t.Error()
	}
	cpu.Execute(inst(decoder.DEX))
	if cpu.core.x != 0 {
		t.Error()
	}

	cpu.Execute(inst(decoder.INY))
	if cpu.core.y != 1 {
		t.Error()
	}
	cpu.Execute(inst(decoder.DEY))
	if cpu.core.y != 0 {
		t.Error()
	}
}

func TestCompare(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.Absolute,
			Value:          0x2000,
		}
	}

	cpu, m := create()
	m.Write(0x2000, 0x80)

	tests := []struct {
		value  uint8
		status status
	}{
		{0x79, status{negative: true}},
		{0x80, status{carry: true, zero: true}},
		{0x81, status{carry: true}},
	}

	for _, test := range tests {
		cpu.core.a = test.value
		cpu.Execute(inst(decoder.CMP))
		if cpu.core.status != test.status {
			t.Error()
		}

		cpu.core.x = test.value
		cpu.Execute(inst(decoder.CPX))
		if cpu.core.status != test.status {
			t.Error()
		}

		cpu.core.y = test.value
		cpu.Execute(inst(decoder.CPY))
		if cpu.core.status != test.status {
			t.Error()
		}
	}
}

func TestCompareFlag(t *testing.T) {
	cpu, _ := create()
	cpu.core.status.overflow = true
	cpu.Execute(decoder.Instruction{
		Op:             decoder.CMP,
		AddressingMode: decoder.Absolute,
		Value:          0x2000,
	})

	if !cpu.core.status.overflow {
		t.Error("clear unexpected flag")
	}
}

func TestCompareNegative(t *testing.T) {
	cpu, _ := create()
	cpu.core.a = 0xff
	cpu.Execute(decoder.Instruction{
		Op:             decoder.CMP,
		AddressingMode: decoder.Absolute,
		Value:          0x2000,
	})

	if !cpu.core.status.negative {
		t.Error("clear unexpected flag")
	}
}

func TestBitCompare(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.Absolute,
			Value:          0x2000,
		}
	}

	cpu, m := create()

	tests := []struct {
		value  uint8
		status status
	}{
		{0x22, status{}},
		{0x0, status{zero: true}},
		{0x80, status{negative: true}},
		{0x70, status{overflow: true}},
	}

	cpu.core.a = 0xff
	for _, test := range tests {
		m.Write(0x2000, test.value)

		i := inst(decoder.BIT)
		cpu.Execute(i)
		if cpu.core.status != test.status {
			t.Errorf("Write %x, then become \n%s, but \n%s",
				test.value,
				test.status.String(),
				cpu.core.status.String())
		}
	}
}

func TestBitFlag(t *testing.T) {
	cpu, _ := create()
	cpu.core.status.carry = true
	cpu.Execute(decoder.Instruction{
		Op:             decoder.BIT,
		AddressingMode: decoder.Absolute,
		Value:          0x2000,
	})

	if !cpu.core.status.carry {
		t.Error("clear unexpected flag")
	}
}

func TestShiftLeft(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.Accumlator,
			Value:          0,
		}
	}

	cpu, _ := create()

	cpu.core.status.carry = true
	cpu.core.a = 0x81
	cpu.Execute(inst(decoder.ASL))
	if cpu.core.a != 0x2 {
		t.Error()
	}
	if cpu.core.status.carry != true {
		t.Error()
	}

	cpu.core.status.carry = true
	cpu.core.a = 0x81
	cpu.Execute(inst(decoder.ROL))
	if cpu.core.a != 0x3 {
		t.Error()
	}
	if cpu.core.status.carry != true {
		t.Error()
	}

	cpu.core.status.carry = true
	cpu.core.a = 0x81
	cpu.Execute(inst(decoder.LSR))
	if cpu.core.a != 0x40 {
		t.Error()
	}
	if cpu.core.status.carry != true {
		t.Error()
	}

	cpu.core.status.carry = true
	cpu.core.a = 0x81
	cpu.Execute(inst(decoder.ROR))
	if cpu.core.a != 0xc0 {
		t.Errorf("%x", cpu.core.a)
	}
	if cpu.core.status.carry != true {
		t.Error()
	}
}

func TestShiftFlag(t *testing.T) {
	cpu, _ := create()

	cpu.core.status.overflow = true

	ops := []decoder.Op{
		decoder.ASL,
		decoder.LSR,
		decoder.ROL,
		decoder.ROR,
	}

	for _, op := range ops {

		cpu.Execute(decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.Accumlator,
			Value:          0,
		})

		if !cpu.core.status.overflow {
			t.Error("clear unexpected flag")
		}
	}
}

func TestPHA(t *testing.T) {
	cpu, m := create()

	cpu.core.s = 0xFF
	cpu.core.a = 0xca

	cpu.Execute(decoder.Instruction{Op: decoder.PHA})

	if m.Read(0x01FF) != 0xca {
		t.Error("cannot push stack")
	}

	if cpu.core.s != 0xFE {
		t.Error("cannot growth stack")
	}

	cpu.core.a = 0
	cpu.Execute(decoder.Instruction{Op: decoder.PLA})

	if cpu.core.a != 0xca {
		t.Error("cannot pop stack")
	}

	if cpu.core.s != 0xFF {
		t.Error("cannot shink stack")
	}
}

func TestPHAFlag(t *testing.T) {
	cpu, m := create()

	tests := []struct {
		value    uint8
		negative bool
		zero     bool
	}{
		{0x00, false, true},
		{0x01, false, false},
		{0x80, true, false},
	}

	for n, test := range tests {
		cpu.core.s = 0x0FE
		m.Write(0x01FF, test.value)
		cpu.core.status.overflow = true
		cpu.Execute(decoder.Instruction{Op: decoder.PLA})

		if cpu.core.status.negative != test.negative {
			t.Errorf("%d cannot store negative flag", n)
		}

		if cpu.core.status.zero != test.zero {
			t.Errorf("%d cannot store zero flag", n)
		}

		if !cpu.core.status.overflow {
			t.Errorf("%d broke overflow flag", n)
		}
	}
}

func TestPHP(t *testing.T) {
	cpu, m := create()

	cpu.core.s = 0xFF
	cpu.core.status = status{negative: true, carry: true}

	cpu.Execute(decoder.Instruction{Op: decoder.PHP})

	if m.Read(0x01FF) != 0x81 {
		t.Error("cannot push stack")
	}

	if cpu.core.s != 0xFE {
		t.Error("cannot growth stack")
	}

	cpu.core.status = status{}
	cpu.Execute(decoder.Instruction{Op: decoder.PLP})

	if (cpu.core.status != status{negative: true, carry: true}) {
		t.Error("cannot pop stack")
	}

	if cpu.core.s != 0xFF {
		t.Error("cannot shink stack")
	}
}

func TestJMP(t *testing.T) {
	cpu, _ := create()
	cpu.Execute(decoder.Instruction{
		Op:             decoder.JMP,
		AddressingMode: decoder.Absolute,
		Value:          0x2000,
	})

	if cpu.core.pc != 0x2000 {
		t.Error("cannot jump expected point")
	}
}

func TestJSR(t *testing.T) {
	cpu, m := create()
	cpu.core.pc = 0xcafe

	cpu.Execute(decoder.Instruction{
		Op:             decoder.JSR,
		AddressingMode: decoder.Absolute,
		Value:          0x3000,
	})

	if cpu.core.pc != 0x3000 {
		t.Error("cannot jump expected point")
	}

	if m.Read16(0x01FE) != 0xcafd {
		t.Errorf("cannot push current pc. %x", m.Read16(0x01FE))
	}

	if cpu.core.s != 0xFD {
		t.Errorf("cannot push down stack: %x", cpu.core.s)
	}

	cpu.Execute(decoder.Instruction{
		Op: decoder.RTS,
	})

	if cpu.core.pc != 0xcafe {
		t.Errorf("cannot return expected point: %x", cpu.core.pc)
	}

	if cpu.core.s != 0xFF {
		t.Error("cannot pop up stack")
	}
}

func TestBranch(t *testing.T) {
	tests := []struct {
		op        decoder.Op
		branch    status
		nonBranch status
	}{
		{decoder.BCS, status{carry: true}, status{}},
		{decoder.BCC, status{}, status{carry: true}},

		{decoder.BEQ, status{zero: true}, status{}},
		{decoder.BNE, status{}, status{zero: true}},

		{decoder.BMI, status{negative: true}, status{}},
		{decoder.BPL, status{}, status{negative: true}},

		{decoder.BVS, status{overflow: true}, status{}},
		{decoder.BVC, status{}, status{overflow: true}},
	}

	cpu, _ := create()
	for n, test := range tests {
		// negative
		cpu.core.pc = 0xcafe
		cpu.core.status = test.branch

		cpu.Execute(decoder.Instruction{
			Op:             test.op,
			AddressingMode: decoder.Relative,
			Value:          0xFF,
		})

		if cpu.core.pc != 0xcafd {
			t.Errorf("%d: cannot branch before point", n)
		}

		// positive
		cpu.core.pc = 0xcafe
		cpu.core.status = test.branch

		cpu.Execute(decoder.Instruction{
			Op:             test.op,
			AddressingMode: decoder.Relative,
			Value:          0x1,
		})

		if cpu.core.pc != 0xcaff {
			t.Errorf("%d: cannot branch after point: %x", n, cpu.core.pc)
		}

		// not jump
		cpu.core.pc = 0xcafe
		cpu.core.status = test.nonBranch
		cpu.Execute(decoder.Instruction{
			Op:             test.op,
			AddressingMode: decoder.Relative,
			Value:          0x80,
		})

		if cpu.core.pc != 0xcafe {
			t.Errorf("%d: unexpected branch happen", n)
		}
	}
}

func TestSetFlag(t *testing.T) {
	tests := []struct {
		op     decoder.Op
		status status
	}{
		{decoder.SEC, status{carry: true}},
		{decoder.SEI, status{irq: true}},
	}

	cpu, _ := create()
	for n, test := range tests {
		cpu.core.status = status{}

		cpu.Execute(decoder.Instruction{
			Op: test.op,
		})

		if cpu.core.status != test.status {
			t.Errorf("%d: unexpected status flag %s", n, cpu.core.status.String())
		}
	}
}

func TestClearFlag(t *testing.T) {
	tests := []struct {
		op     decoder.Op
		status status
	}{
		{decoder.CLC, status{carry: true}},
		{decoder.CLI, status{irq: true}},
		{decoder.CLV, status{overflow: true}},
	}

	cpu, _ := create()
	for n, test := range tests {
		cpu.core.status = test.status

		cpu.Execute(decoder.Instruction{
			Op: test.op,
		})

		if (cpu.core.status != status{}) {
			t.Errorf("%d: clear status flag", n)
		}
	}
}

func TestBrk(t *testing.T) {
	cpu, _ := create()

	cpu.Execute(decoder.Instruction{
		Op: decoder.BRK,
	})

	if cpu.interrupt != brk {
		t.Error()
	}
}

func TestRti(t *testing.T) {
	cpu, _ := create()

	cpu.core.push(0xca)
	cpu.core.push(0xfe)
	cpu.core.push(0x80)

	cpu.Execute(decoder.Instruction{
		Op: decoder.RTI,
	})

	if cpu.core.pc != 0xcafe {
		t.Errorf("%x", cpu.core.pc)
	}

	if !cpu.core.status.negative {
		t.Error()
	}
}
