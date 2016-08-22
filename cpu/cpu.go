package cpu

import (
	d "github.com/mzp/famicom/decoder"
	"github.com/mzp/famicom/memory"
)

type CPU struct {
	core      *core
	interrupt interrupt
}

func New(m *memory.Memory, start int) *CPU {
	core := makeCore(m, start)
	t := CPU{core: core}
	return &t
}

func (self *CPU) decode() (d.Instruction, int) {
	return d.Decode(self.core.memory.Data[:], self.core.pc)
}

func (self *CPU) CurrentInstruction() d.Instruction {
	inst, _ := self.decode()
	return inst
}

func (self *CPU) InterrruptReset() {
	self.interrupt = reset
}

func (self *CPU) InterrruptNMI() {
	self.interrupt = nmi
}

func (self *CPU) InterrruptIrq() {
	self.interrupt = irq
}

func (self *CPU) InterrruptBreak() {
	self.interrupt = brk
}

func (self *CPU) Fetch() d.Instruction {
	inst, n := self.decode()
	self.core.pc += n
	return inst
}

func (c *CPU) Execute(inst d.Instruction) {
	if inst.Op == d.BRK {
		c.interrupt = brk
	} else {
		execute(c.core, inst)
	}
}

func (c *CPU) Step() {
	switch c.interrupt {
	case reset:
		handleReset(c.core)
	case nmi:
		handleNMI(c.core)
	case irq:
		if c.core.status.irq {
			c.Execute(c.Fetch())
		} else {
			handleInteruption(c.core, false)
		}
	case brk:
		if c.core.status.irq {
			c.Execute(c.Fetch())
		} else {
			handleInteruption(c.core, true)
		}
	case none:
		c.Execute(c.Fetch())
	default:
		panic("must not happen")
	}
	c.interrupt = none
}

func (c *CPU) String() string {
	return c.core.String()
}
