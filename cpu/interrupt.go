package cpu

type interrupt int

const (
	none interrupt = iota
	reset
	nmi
	irq
	brk
)

func handleReset(c *core) {
	c.pc = int(c.memory.Read16(0xFFFC))
	c.status.irq = true
}

func handleNMI(c *core) {
	c.push(uint8(c.pc >> 8))
	c.push(uint8(c.pc))
	c.push(c.status.uint8())
	c.status.irq = true
	c.status.brk = false
	c.pc = int(c.memory.Read16(0xFFFA))
}

func handleInteruption(c *core, brk bool) {
	c.push(uint8(c.pc >> 8))
	c.push(uint8(c.pc))
	c.push(c.status.uint8())
	c.status.irq = true
	c.status.brk = brk
	c.pc = int(c.memory.Read16(0xFFFE))
}
