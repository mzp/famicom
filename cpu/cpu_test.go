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

func TestAddressingMode(t *testing.T) {
	cpu, m := create()
	m.Write(0x01, 0x00)
	m.Write(0x02, 0x20)

	m.Write(0x80, 0x42)
	m.Write(0x81, 0x43)
	m.Write(0x82, 0x44)

	m.Write(0x2000, 0x80)
	m.Write(0x2001, 0x81)
	m.Write(0x2002, 0x82)

	cpu.x = 1
	cpu.y = 2

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
			Op:             decoder.LDA,
			AddressingMode: test.mode,
			Value:          test.value,
		}
		cpu.Execute(inst)

		if cpu.a != test.expect {
			t.Errorf("%s %x != %x", inst.String(), cpu.a, test.expect)
		}
	}
}

func TestNZ(t *testing.T) {
	lda := func(value int) decoder.Instruction {
		return decoder.Instruction{
			Op:             decoder.LDA,
			AddressingMode: decoder.Immediate,
			Value:          value,
		}
	}

	cpu, _ := create()

	tests := []struct {
		value  int
		status status
	}{
		{0, status{zero: true}},
		{1, status{}},
		{0x80, status{negative: true}},
	}

	for _, test := range tests {
		cpu.Execute(lda(test.value))
		if cpu.status != test.status {
			t.Errorf("load %x, %s must be %s", test.value, cpu.status.String(), test.status.String())
		}
	}
}

func TestLoad(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.Immediate,
			Value:          42,
		}
	}

	cpu, _ := create()

	cpu.Execute(inst(decoder.LDA))
	if cpu.a != 42 {
		t.Errorf("cpu.a = %x, but must be 42", cpu.a)
	}

	cpu.Execute(inst(decoder.LDX))
	if cpu.x != 42 {
		t.Errorf("cpu.x = %x, but must be 42", cpu.a)
	}

	cpu.Execute(inst(decoder.LDY))
	if cpu.y != 42 {
		t.Errorf("cpu.y = %x, but must be 42", cpu.a)
	}
}

func TestStore(t *testing.T) {
	inst := func(op decoder.Op) decoder.Instruction {
		return decoder.Instruction{
			Op:             op,
			AddressingMode: decoder.Absolute,
			Value:          42,
		}
	}

	cpu, m := create()

	cpu.a = 1
	cpu.Execute(inst(decoder.STA))
	if m.Read(42) != 1 {
		t.Errorf("Store 1, but %x", m.Read(42))
	}

	cpu.x = 2
	cpu.Execute(inst(decoder.STX))
	if m.Read(42) != 2 {
		t.Errorf("Store 2, but %x", m.Read(42))
	}

	cpu.y = 3
	cpu.Execute(inst(decoder.STY))
	if m.Read(42) != 3 {
		t.Errorf("Store 3, but %x", m.Read(42))
	}
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

	cpu.a = 1
	cpu.Execute(inst(decoder.TAX))
	if cpu.x != 1 {
		t.Errorf("Transfer 1, but %x", cpu.x)
	}

	cpu.Execute(inst(decoder.TAY))
	if cpu.y != 1 {
		t.Errorf("Transfer 1, but %x", cpu.y)
	}

	cpu.s = 2
	cpu.Execute(inst(decoder.TSX))
	if cpu.x != 2 {
		t.Errorf("Transfer 2, but %x", cpu.x)
	}

	cpu.x = 3
	cpu.Execute(inst(decoder.TXA))
	if cpu.a != 3 {
		t.Errorf("Transfer 3, but %x", cpu.a)
	}

	cpu.Execute(inst(decoder.TXS))
	if cpu.s != 3 {
		t.Errorf("Transfer 3, but %x", cpu.s)
	}

	cpu.y = 4
	cpu.Execute(inst(decoder.TYA))
	if cpu.a != 4 {
		t.Errorf("Transfer 4, but %x", cpu.a)
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

		// TODO
		// {0x80, 0x80, false, 0x00, status{overflow: true, carry: true}},
	}

	for _, test := range tests {
		cpu.a = test.a
		cpu.status.carry = test.carry
		m.Write(0x2000, test.memory)

		cpu.Execute(inst)

		if cpu.a != test.expect {
			t.Errorf("Expect %x, but %x", test.expect, cpu.a)
		}

		if cpu.status != test.status {
			t.Errorf("Expect %s, but %s", test.status.String(), cpu.status.String())
		}
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
		{0x30, 0x20, true, 0x10, status{}},
		{0x0, 0x0, true, 0x0, status{zero: true}},
		{0x2, 0x0, false, 0x1, status{}},
		{0x0, 0x1, true, 0xFF, status{negative: true}},
		// TODO:
		// {0x80, 0x1, true, 0x7F, status{overflow: true}},
	}

	for _, test := range tests {
		cpu.a = test.a
		cpu.status.carry = test.carry
		m.Write(0x2000, test.memory)

		cpu.Execute(inst)

		if cpu.a != test.expect {
			t.Errorf("Expect %x, but %x", test.expect, cpu.a)
		}

		if cpu.status != test.status {
			t.Errorf("Expect %s, but %s", test.status.String(), cpu.status.String())
		}
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
	cpu.a = 0xad
	cpu.Execute(inst(decoder.AND))
	if cpu.a != 0x8c {
		t.Errorf("0xde & 0xad = %x", cpu.a)
	}

	cpu.a = 0xad
	cpu.Execute(inst(decoder.EOR))
	if cpu.a != 0x73 {
		t.Errorf("0xde ^ 0xad = %x", cpu.a)
	}

	cpu.a = 0xad
	cpu.Execute(inst(decoder.ORA))
	if cpu.a != 0xff {
		t.Errorf("0xde | 0xad = %x", cpu.a)
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
	if cpu.x != 1 {
		t.Error()
	}
	cpu.Execute(inst(decoder.DEX))
	if cpu.x != 0 {
		t.Error()
	}

	cpu.Execute(inst(decoder.INY))
	if cpu.y != 1 {
		t.Error()
	}
	cpu.Execute(inst(decoder.DEY))
	if cpu.y != 0 {
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
		cpu.a = test.value
		cpu.Execute(inst(decoder.CMP))
		if cpu.status != test.status {
			t.Error()
		}

		cpu.x = test.value
		cpu.Execute(inst(decoder.CPX))
		if cpu.status != test.status {
			t.Error()
		}

		cpu.y = test.value
		cpu.Execute(inst(decoder.CPY))
		if cpu.status != test.status {
			t.Error()
		}
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

	cpu.a = 0xff
	for _, test := range tests {
		m.Write(0x2000, test.value)

		i := inst(decoder.BIT)
		cpu.Execute(i)
		if cpu.status != test.status {
			t.Errorf("Write %x, then become \n%s, but \n%s",
				test.value,
				test.status.String(),
				cpu.status.String())
		}
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

	cpu.status.carry = true
	cpu.a = 0x81
	cpu.Execute(inst(decoder.ASL))
	if cpu.a != 0x2 {
		t.Error()
	}
	if cpu.status.carry != true {
		t.Error()
	}

	cpu.status.carry = true
	cpu.a = 0x81
	cpu.Execute(inst(decoder.ROL))
	if cpu.a != 0x3 {
		t.Error()
	}
	if cpu.status.carry != true {
		t.Error()
	}

	cpu.status.carry = true
	cpu.a = 0x81
	cpu.Execute(inst(decoder.LSR))
	if cpu.a != 0x40 {
		t.Error()
	}
	if cpu.status.carry != true {
		t.Error()
	}

	cpu.status.carry = true
	cpu.a = 0x81
	cpu.Execute(inst(decoder.ROR))
	if cpu.a != 0xc0 {
		t.Errorf("%x", cpu.a)
	}
	if cpu.status.carry != true {
		t.Error()
	}
}