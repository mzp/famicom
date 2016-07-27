package cpu

import "fmt"

type status struct {
	negative, overflow, brk, irq, zero, carry bool
}

func toInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func (status status) String() string {
	return fmt.Sprintf("N:%d V:%d B%d I:%d Z:%d C:%d",
		toInt(status.negative),
		toInt(status.overflow),
		toInt(status.brk),
		toInt(status.irq),
		toInt(status.zero),
		toInt(status.carry))
}
