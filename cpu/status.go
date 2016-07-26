package cpu

import "fmt"

type status struct {
	negative, overflow, brk, irq, zero, carry bool
}

func (status status) String() string {
	b := func(x bool) int {
		if x {
			return 1
		} else {
			return 0
		}
	}

	return fmt.Sprintf("N:%d V:%d B%d I:%d Z:%d C:%d",
		b(status.negative),
		b(status.overflow),
		b(status.brk),
		b(status.irq),
		b(status.zero),
		b(status.carry))
}
