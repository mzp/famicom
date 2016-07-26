package cpu

import (
	"fmt"

	d "github.com/mzp/famicom/decoder"
	memlib "github.com/mzp/famicom/memory"
)

type status struct {
	negative, overflow, brk, irq, zero, carry bool
}

type CPU struct {
	memory  *memlib.Memory
	pc      int
	a, x, y byte
	status  status
}

func New(m *memlib.Memory, start int) *CPU {
	t := CPU{memory: m, pc: start}
	return &t
}

func (c *CPU) CurrentInstruction() d.Instruction {
	inst, _ := d.Decode(c.memory.Data[:], c.pc)
	return inst
}

func (c *CPU) address(inst d.Instruction) uint16 {
	value := inst.Value

	switch inst.AddressingMode {
	case d.ZeroPage:
		return uint16(value)
	case d.ZeroPageX:
		return uint16(value) + uint16(c.x)
	case d.ZeroPageY:
		return uint16(value) + uint16(c.y)
	case d.Absolute:
		return uint16(value)
	case d.AbsoluteX:
		return uint16(value) + uint16(c.x)
	case d.AbsoluteY:
		return uint16(value) + uint16(c.y)
	case d.Indirect:
		return c.memory.Read16(uint16(value))
	case d.IndirectX:
		return c.memory.Read16(uint16(value) + uint16(c.x))
	case d.IndirectY:
		address := c.memory.Read16(uint16(value))
		return address + uint16(c.y)
	default:
		panic("unknown addressing mode")
	}
}

func (c *CPU) read(inst d.Instruction) uint8 {
	value := inst.Value

	switch inst.AddressingMode {
	case d.Immediate:
		return uint8(value)
	default:
		address := c.address(inst)
		return c.memory.Read(address)
	}
}

func nz(value uint8) status {
	return status{
		negative: (0x80 & value) != 0,
		zero:     value == 0,
	}
}

func (c *CPU) setA(value uint8) {
	c.a = value
	c.status = nz(value)
}

func (c *CPU) setX(value uint8) {
	c.x = value
	c.status = nz(value)
}

func (c *CPU) setY(value uint8) {
	c.y = value
	c.status = nz(value)
}

func (c *CPU) clear() {
	c.status = status{}
}

func (c *CPU) Step() {
	inst, n := d.Decode(c.memory.Data[:], c.pc)
	c.pc += n

	switch inst.Op {
	case d.LDA:
		c.setA(c.read(inst))
	case d.LDX:
		c.setX(c.read(inst))
	case d.LDY:
		c.setY(c.read(inst))
	default:
		c.status = status{}
	}
}

func (c *CPU) String() string {
	return fmt.Sprintf("x:%08x y:%08x a:%08x", c.x, c.y, c.a)
}
