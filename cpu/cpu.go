package cpu

import (
	"fmt"

	"github.com/mzp/famicom/bits"
	d "github.com/mzp/famicom/decoder"
	memlib "github.com/mzp/famicom/memory"
)

type interrupt int

const (
	none interrupt = iota
	reset
	nmi
	irq
	brk
)

type CPU struct {
	memory     *memlib.Memory
	pc         int
	a, x, y, s byte
	status     status
	interrupt  interrupt
}

func New(m *memlib.Memory, start int) *CPU {
	t := CPU{memory: m, pc: start, s: 0xff}
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
	case d.Relative:
		return uint16(c.pc + int(int8(value)))
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

func (cpu *CPU) nz(value uint8) {
	cpu.status.negative = bits.IsFlag(value, 7)
	cpu.status.zero = value == 0
}

func (c *CPU) load(reg *uint8, value uint8) {
	*reg = value
	c.nz(value)
}

func (c *CPU) store(address uint16, value uint8) {
	c.memory.Write(address, value)
}

func (c *CPU) adc(inst d.Instruction) {
	value := int(c.a) + int(c.read(inst)) + toInt(c.status.carry)
	c.nz(uint8(value))
	c.status.carry = value > 0xFF
	c.status.overflow =
		(c.a <= 0x7F && 0x80 <= uint8(value)) ||
			(uint8(value) <= 0x7F && 0x80 <= c.a)
	c.a = uint8(value)
}

func (c *CPU) sbc(inst d.Instruction) {
	value := int(c.a) - int(c.read(inst)) - (1 - toInt(c.status.carry))

	c.nz(uint8(value))
	c.status.carry = value >= 0
	c.status.overflow =
		(c.a <= 0x7F && 0x80 <= uint8(value)) ||
			(uint8(value) <= 0x7F && 0x80 <= c.a)

	c.a = uint8(value)
}

func (c *CPU) compare(reg uint8, inst d.Instruction) {
	value := c.read(inst)

	c.status.negative = reg < value
	c.status.zero = reg == value
	c.status.carry = reg >= value
}

func (c *CPU) shiftL(inst d.Instruction, carry bool) {
	bit := uint8(toInt(carry))
	if inst.AddressingMode == d.Accumlator {
		ret := c.a<<1 | bit
		c.status.negative = bits.IsFlag(ret, 7)
		c.status.zero = ret == 0
		c.status.carry = bits.IsFlag(c.a, 7)
		c.a = ret
	} else {
		address := c.address(inst)
		value := c.memory.Read(address)
		ret := value<<1 | bit
		c.status.negative = bits.IsFlag(ret, 7)
		c.status.zero = ret == 0
		c.status.carry = bits.IsFlag(value, 7)
		c.memory.Write(address, ret)
	}
}

func (c *CPU) shiftR(inst d.Instruction, carry bool) {
	var bit uint8

	if carry {
		bit = 0x80
	}

	if inst.AddressingMode == d.Accumlator {
		ret := c.a>>1 | bit
		c.status.negative = bits.IsFlag(ret, 7)
		c.status.zero = ret == 0
		c.status.carry = bits.IsFlag(c.a, 0)
		c.a = ret
	} else {
		address := c.address(inst)
		value := c.memory.Read(address)
		ret := value>>1 | bit
		c.status.negative = bits.IsFlag(ret, 7)
		c.status.zero = ret == 0
		c.status.carry = bits.IsFlag(value, 0)
		c.memory.Write(address, ret)
	}
}

func (c *CPU) push(value uint8) {
	c.memory.Write(uint16(c.s)+0x100, value)
	c.s -= 1
}

func (c *CPU) InterrruptReset() {
	c.interrupt = reset
}

func (c *CPU) InterrruptNMI() {
	c.interrupt = nmi
}

func (c *CPU) InterrruptIrq() {
	c.interrupt = irq
}

func (c *CPU) InterrruptBreak() {
	c.interrupt = brk
}

func (c *CPU) Fetch() d.Instruction {
	inst, n := d.Decode(c.memory.Data[:], c.pc)
	c.pc += n
	return inst
}

func (c *CPU) Execute(inst d.Instruction) {
	switch inst.Op {
	case d.LDA:
		c.load(&c.a, c.read(inst))
	case d.LDX:
		c.load(&c.x, c.read(inst))
	case d.LDY:
		c.load(&c.y, c.read(inst))
	case d.STA:
		c.store(c.address(inst), c.a)
	case d.STX:
		c.store(c.address(inst), c.x)
	case d.STY:
		c.store(c.address(inst), c.y)
	case d.TAX:
		c.load(&c.x, c.a)
	case d.TAY:
		c.load(&c.y, c.a)
	case d.TSX:
		c.load(&c.x, c.s)
	case d.TXA:
		c.load(&c.a, c.x)
	case d.TXS:
		c.load(&c.s, c.x)
	case d.TYA:
		c.load(&c.a, c.y)
	case d.ADC:
		c.adc(inst)
	case d.AND:
		c.load(&c.a, c.a&c.read(inst))
	case d.ASL:
		c.shiftL(inst, false)
	case d.BIT:
		value := c.read(inst)
		c.status.negative = bits.IsFlag(value, 7)
		c.status.overflow = bits.IsFlag(value, 6)
		c.status.zero = c.a&value == 0
	case d.CMP:
		c.compare(c.a, inst)
	case d.CPX:
		c.compare(c.x, inst)
	case d.CPY:
		c.compare(c.y, inst)
	case d.DEC:
		value := c.read(inst) - 1
		c.nz(value)
		c.store(c.address(inst), value)
	case d.DEX:
		c.load(&c.x, c.x-1)
	case d.DEY:
		c.load(&c.y, c.y-1)
	case d.EOR:
		value := c.read(inst)
		c.load(&c.a, c.a^value)
	case d.INC:
		value := c.read(inst)
		c.store(c.address(inst), value+1)
	case d.INX:
		c.load(&c.x, c.x+1)
	case d.INY:
		c.load(&c.y, c.y+1)
	case d.LSR:
		c.shiftR(inst, false)
	case d.ROL:
		c.shiftL(inst, c.status.carry)
	case d.ROR:
		c.shiftR(inst, c.status.carry)
	case d.ORA:
		value := c.read(inst)
		c.load(&c.a, c.a|value)
	case d.SBC:
		c.sbc(inst)
	case d.PHA:
		c.push(c.a)
	case d.PLA:
		c.s += 1
		c.a = c.memory.Read(uint16(c.s) + 0x100)
		c.nz(c.a)
	case d.PHP:
		c.push(c.status.uint8())
	case d.PLP:
		c.s += 1
		value := c.memory.Read(uint16(c.s) + 0x100)
		c.status.set(value)
	case d.JMP:
		c.pc = int(c.address(inst))
	case d.JSR:
		pc := c.pc - 1
		c.push(uint8(pc >> 8))
		c.push(uint8(pc))
		c.pc = int(c.address(inst))
	case d.RTS:
		c.s += 1
		value := c.memory.Read16(uint16(c.s) + 0x100)
		c.s += 1
		c.pc = int(value) + 1
	case d.RTI:
		c.s += 1
		status := c.memory.Read(uint16(c.s) + 0x100)
		c.s += 1
		pc := c.memory.Read16(uint16(c.s) + 0x100)
		c.s += 1
		c.pc = int(pc)
		c.status.set(status)
	case d.BCS:
		if c.status.carry {
			c.pc = int(c.address(inst))
		}
	case d.BCC:
		if !c.status.carry {
			c.pc = int(c.address(inst))
		}
	case d.BEQ:
		if c.status.zero {
			c.pc = int(c.address(inst))
		}
	case d.BNE:
		if !c.status.zero {
			c.pc = int(c.address(inst))
		}
	case d.BMI:
		if c.status.negative {
			c.pc = int(c.address(inst))
		}
	case d.BPL:
		if !c.status.negative {
			c.pc = int(c.address(inst))
		}
	case d.BVS:
		if c.status.overflow {
			c.pc = int(c.address(inst))
		}
	case d.BVC:
		if !c.status.overflow {
			c.pc = int(c.address(inst))
		}
	case d.SEC:
		c.status.carry = true
	case d.SEI:
		c.status.irq = true
	case d.CLC:
		c.status.carry = false
	case d.CLI:
		c.status.irq = false
	case d.CLV:
		c.status.overflow = false
	case d.BRK:
		c.interrupt = brk
	case d.NOP:
		// nothing to do
	default:
	}
}

func (c *CPU) Step() {
	switch c.interrupt {
	case reset:
		c.pc = int(c.memory.Read16(0xFFFC))
		c.status.irq = true
	case nmi:
		c.push(uint8(c.pc >> 8))
		c.push(uint8(c.pc))
		c.push(c.status.uint8())
		c.status.irq = true
		c.status.brk = false
		c.pc = int(c.memory.Read16(0xFFFA))
	case irq:
		if c.status.irq {
			c.Execute(c.Fetch())
		} else {
			c.push(uint8(c.pc >> 8))
			c.push(uint8(c.pc))
			c.push(c.status.uint8())
			c.status.irq = true
			c.status.brk = false
			c.pc = int(c.memory.Read16(0xFFFE))
		}
	case brk:
		if c.status.irq {
			c.Execute(c.Fetch())
		} else {
			c.push(uint8(c.pc >> 8))
			c.push(uint8(c.pc))
			c.push(c.status.uint8())
			c.status.irq = true
			c.status.brk = true
			c.pc = int(c.memory.Read16(0xFFFE))
		}
	case none:
		c.Execute(c.Fetch())
	default:
		panic("must not happen")
	}
	c.interrupt = none
}

func (c *CPU) String() string {
	return fmt.Sprintf("x:%08x y:%08x a:%08x s:%08x [%s]", c.x, c.y, c.a, c.s, c.status.String())
}
