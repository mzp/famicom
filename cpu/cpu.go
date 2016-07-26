package cpu

import (
	"fmt"

	"github.com/mzp/famicom/decoder"
	memlib "github.com/mzp/famicom/memory"
)

type CPU struct {
	memory  *memlib.Memory
	pc      int
	a, x, y byte
}

func New(m *memlib.Memory, start int) *CPU {
	t := CPU{memory: m, pc: start}
	return &t
}

func (c *CPU) CurrentInstruction() decoder.Instruction {
	inst, _ := decoder.Decode(c.memory.Data[:], c.pc)
	return inst
}

func (c *CPU) read(inst decoder.Instruction) uint8 {
	value := inst.Value

	switch inst.AddressingMode {
	case decoder.Immediate:
		return uint8(value)
	case decoder.ZeroPage:
		return c.memory.ReadZeroPage(uint8(value))
	case decoder.ZeroPageX:
		return c.memory.ReadZeroPage(uint8(value) + c.x)
	case decoder.ZeroPageY:
		return c.memory.ReadZeroPage(uint8(value) + c.y)
	case decoder.Absolute:
		return c.memory.Read(uint16(value))
	case decoder.AbsoluteX:
		return c.memory.Read(uint16(value) + uint16(c.x))
	case decoder.AbsoluteY:
		return c.memory.Read(uint16(value) + uint16(c.y))
	case decoder.Indirect:
		address := c.memory.Read16(uint16(value))
		return c.memory.Read(address)
	case decoder.IndirectX:
		address := c.memory.Read16(uint16(value) + uint16(c.x))
		return c.memory.Read(address)
	case decoder.IndirectY:
		address := c.memory.Read16(uint16(value))
		return c.memory.Read(address + uint16(c.y))
	default:
		panic("unknown addressing mode")
	}
}

func (c *CPU) Step() {
	inst, n := decoder.Decode(c.memory.Data[:], c.pc)
	c.pc += n

	switch inst.Op {
	case decoder.LDA:
		c.a = c.read(inst)
	case decoder.LDX:
		c.x = c.read(inst)
	case decoder.LDY:
		c.y = c.read(inst)
	}
}

func (c *CPU) String() string {
	return fmt.Sprintf("x:%08x y:%08x a:%08x", c.x, c.y, c.a)
}
