package cpu

import (
	"fmt"

	"github.com/mzp/famicom/bits"
	"github.com/mzp/famicom/decoder"
	"github.com/mzp/famicom/memory"
)

type core struct {
	memory     *memory.Memory
	pc         int
	a, x, y, s byte
	status     status
}

func makeCore(m *memory.Memory, start int) *core {
	t := core{memory: m, pc: start, s: 0xff}
	return &t
}

func (self *core) address(inst decoder.Instruction) uint16 {
	value := inst.Value

	switch inst.AddressingMode {
	case decoder.ZeroPage:
		return uint16(value)
	case decoder.ZeroPageX:
		return uint16(value) + uint16(self.x)
	case decoder.ZeroPageY:
		return uint16(value) + uint16(self.y)
	case decoder.Absolute:
		return uint16(value)
	case decoder.AbsoluteX:
		return uint16(value) + uint16(self.x)
	case decoder.AbsoluteY:
		return uint16(value) + uint16(self.y)
	case decoder.Indirect:
		return self.memory.Read16(uint16(value))
	case decoder.IndirectX:
		return self.memory.Read16(uint16(value) + uint16(self.x))
	case decoder.IndirectY:
		address := self.memory.Read16(uint16(value))
		return address + uint16(self.y)
	case decoder.Relative:
		return uint16(self.pc + int(int8(value)))
	default:
		panic("unknown addressing mode")
	}
}

func (self *core) read(inst decoder.Instruction) uint8 {
	value := inst.Value

	switch inst.AddressingMode {
	case decoder.Immediate:
		return uint8(value)
	default:
		address := self.address(inst)
		return self.memory.Read(address)
	}
}

func (self *core) nz(value uint8) {
	self.status.negative = bits.IsFlag(value, 7)
	self.status.zero = value == 0
}

func (self *core) load(reg *uint8, value uint8) {
	*reg = value
	self.nz(value)
}

func (self *core) store(address uint16, value uint8) {
	self.memory.Write(address, value)
}

func (self *core) push(value uint8) {
	self.memory.Write(uint16(self.s)+0x100, value)
	self.s -= 1
}

func (self *core) push16(value int) {
	self.push(uint8(value >> 8))
	self.push(uint8(value))
}

func (self *core) pop() uint8 {
	self.s += 1
	return self.memory.Read(uint16(self.s) + 0x100)
}

func (self *core) pop16() uint16 {
	self.s += 1
	value := self.memory.Read16(uint16(self.s) + 0x100)
	self.s += 1
	return value
}

func (self *core) String() string {
	return fmt.Sprintf("x:%08x y:%08x a:%08x s:%08x [%s]", self.x, self.y, self.a, self.s, self.status.String())
}
