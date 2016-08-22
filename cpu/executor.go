package cpu

import (
	"github.com/mzp/famicom/bits"
	d "github.com/mzp/famicom/decoder"
)

func execute(c *core, inst d.Instruction) {
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
		adc(c, inst)
	case d.AND:
		c.load(&c.a, c.a&c.read(inst))
	case d.ASL:
		shiftL(c, inst, false)
	case d.BIT:
		value := c.read(inst)
		c.status.negative = bits.IsFlag(value, 7)
		c.status.overflow = bits.IsFlag(value, 6)
		c.status.zero = c.a&value == 0
	case d.CMP:
		compare(c, c.a, inst)
	case d.CPX:
		compare(c, c.x, inst)
	case d.CPY:
		compare(c, c.y, inst)
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
		value := c.read(inst) + 1
		c.nz(value)
		c.store(c.address(inst), value)
	case d.INX:
		c.load(&c.x, c.x+1)
	case d.INY:
		c.load(&c.y, c.y+1)
	case d.LSR:
		shiftR(c, inst, false)
	case d.ROL:
		shiftL(c, inst, c.status.carry)
	case d.ROR:
		shiftR(c, inst, c.status.carry)
	case d.ORA:
		value := c.read(inst)
		c.load(&c.a, c.a|value)
	case d.SBC:
		sbc(c, inst)
	case d.PHA:
		c.push(c.a)
	case d.PLA:
		c.a = c.pop()
		c.nz(c.a)
	case d.PHP:
		c.push(c.status.uint8())
	case d.PLP:
		value := c.pop()
		c.status.set(value)
	case d.JMP:
		c.pc = int(c.address(inst))
	case d.JSR:
		pc := c.pc - 1
		c.push16(pc)
		c.pc = int(c.address(inst))
	case d.RTS:
		value := c.pop16()
		c.pc = int(value) + 1
	case d.RTI:
		status := c.pop()
		pc := c.pop16()
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
	case d.NOP:
		// nothing to do
	default:
	}
}

func adc(c *core, inst d.Instruction) {
	value := int(c.a) + int(c.read(inst)) + toInt(c.status.carry)
	c.nz(uint8(value))
	c.status.carry = value > 0xFF
	c.status.overflow =
		(c.a <= 0x7F && 0x80 <= uint8(value)) ||
			(uint8(value) <= 0x7F && 0x80 <= c.a)
	c.a = uint8(value)
}

func sbc(c *core, inst d.Instruction) {
	value := int(c.a) - int(c.read(inst)) - (1 - toInt(c.status.carry))

	c.nz(uint8(value))
	c.status.carry = value >= 0
	c.status.overflow =
		(c.a <= 0x7F && 0x80 <= uint8(value)) ||
			(uint8(value) <= 0x7F && 0x80 <= c.a)

	c.a = uint8(value)
}

func compare(c *core, reg uint8, inst d.Instruction) {
	value := c.read(inst)
	c.nz(reg - value)
	c.status.carry = reg >= value
}

func shiftL(c *core, inst d.Instruction, carry bool) {
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

func shiftR(c *core, inst d.Instruction, carry bool) {
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
