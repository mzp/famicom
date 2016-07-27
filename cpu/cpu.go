package cpu

import (
	"fmt"

	d "github.com/mzp/famicom/decoder"
	memlib "github.com/mzp/famicom/memory"
)

type CPU struct {
	memory     *memlib.Memory
	pc         int
	a, x, y, s byte
	status     status
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

func nzc(value int) status {
	return status{
		negative: (0x80 & value) != 0,
		zero:     value == 0,
		carry:    value > 0xFF,
	}
}

func nvzc(value int) status {
	return status{
		negative: (0x80 & value) != 0,
		overflow: value > 0x7F,
		zero:     value == 0,
		carry:    value > 0xFF,
	}
}

func (c *CPU) load(reg *uint8, value uint8) {
	*reg = value
	c.status = nz(value)
}

func (c *CPU) store(address uint16, value uint8) {
	c.memory.Write(address, value)
	c.status = status{}
}

func (c *CPU) loadVC(reg *uint8, value int) {
	*reg = uint8(value)
	c.status = nvzc(value)
}

func (c *CPU) loadC(reg *uint8, value int) {
	*reg = uint8(value)
	c.status = nzc(value)
}

func (c *CPU) storeC(address uint16, value int) {
	c.memory.Write(address, uint8(value))
	c.status = nzc(value)
}

func (c *CPU) compare(reg uint8, inst d.Instruction) {
	value := c.read(inst)

	c.status = status{
		negative: reg < value,
		zero:     reg == value,
		carry:    reg >= value,
	}
}

func (c *CPU) shiftL(inst d.Instruction, carry bool) {
	bit := uint8(toInt(carry))
	if inst.AddressingMode == d.Accumlator {
		ret := c.a<<1 | bit
		c.status = status{
			negative: (0x80 & ret) != 0,
			zero:     ret == 0,
			carry:    (c.a & 0x80) != 0,
		}
		c.a = ret
	} else {
		address := c.address(inst)
		value := c.memory.Read(address)
		ret := value<<1 | bit
		c.status = status{
			negative: (0x80 & ret) != 0,
			zero:     ret == 0,
			carry:    (value & 0x80) != 0,
		}
		c.memory.Write(address, ret)
	}
}

func (c *CPU) shiftR(inst d.Instruction, carry bool) {
	var bit uint8

	if inst.AddressingMode == d.Accumlator {
		ret := c.a>>1 | bit
		c.status = status{
			negative: (0x80 & ret) != 0,
			zero:     ret == 0,
			carry:    (c.a & 0x01) != 0,
		}
		c.a = ret
	} else {
		address := c.address(inst)
		value := c.memory.Read(address)
		ret := value>>1 | bit
		c.status = status{
			negative: (0x80 & ret) != 0,
			zero:     ret == 0,
			carry:    (value & 0x01) != 0,
		}
		c.memory.Write(address, ret)
	}
}

func (c *CPU) Step() {
	inst, n := d.Decode(c.memory.Data[:], c.pc)
	c.pc += n

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
		c.loadVC(&c.a,
			int(c.a)+int(c.read(inst))+toInt(c.status.carry))
	case d.AND:
		c.load(&c.a, c.a&c.read(inst))
	case d.ASL:
		c.shiftL(inst, false)
	case d.BIT:
		value := c.read(inst)
		c.status = status{
			negative: value&0x80 != 0,
			overflow: value&0x40 != 0,
			zero:     c.a&value == 0,
		}
	case d.CMP:
		c.compare(c.a, inst)
	case d.CPX:
		c.compare(c.x, inst)
	case d.CPY:
		c.compare(c.y, inst)
	case d.DEC:
		value := c.read(inst)
		c.store(c.address(inst), value-1)
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
		c.loadVC(&c.a,
			int(c.a)-int(c.read(inst))-(1-toInt(c.status.carry)))
	default:
		c.status = status{}
	}
}

func (c *CPU) String() string {
	return fmt.Sprintf("x:%08x y:%08x a:%08x s:%08x [%s]", c.x, c.y, c.a, c.s, c.status.String())
}
