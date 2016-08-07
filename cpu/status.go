package cpu

import (
	"fmt"

	"github.com/mzp/famicom/bits"
)

type status struct {
	negative, overflow, brk, irq, zero, carry bool
}

func (s *status) set(flag uint8) {
	s.negative = bits.IsFlag(flag, 7)
	s.overflow = bits.IsFlag(flag, 6)
	s.brk = bits.IsFlag(flag, 4)
	s.irq = bits.IsFlag(flag, 2)
	s.zero = bits.IsFlag(flag, 1)
	s.carry = bits.IsFlag(flag, 0)
}

func (s *status) uint8() uint8 {
	var flag uint8

	if s.negative {
		flag |= 0x80
	}
	if s.overflow {
		flag |= 0x40
	}
	if s.brk {
		flag |= 0x10
	}
	if s.irq {
		flag |= 0x04
	}
	if s.zero {
		flag |= 0x02
	}
	if s.carry {
		flag |= 0x01
	}

	return flag
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
