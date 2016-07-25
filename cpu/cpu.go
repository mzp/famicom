package cpu

import (
	memlib "github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/decoder"
)

type CPU struct {
	memory *memlib.Memory
	pc int
}

func New(m *memlib.Memory, start int) *CPU {
	t := CPU{ m, start }
	return &t
}

func (c *CPU) CurrentInstruction() decoder.Instruction {
	inst, _ := decoder.Decode(c.memory.Data[:], c.pc)
	return inst
}

func (c *CPU) Step() {
	_, n := decoder.Decode(c.memory.Data[:], c.pc)
	c.pc += n
}
