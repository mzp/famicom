package cpu

import "fmt"

type status struct {
	negative, overflow, brk, irq, zero, carry bool
}

func (s *status) set(flag uint8) {
	s.negative = (flag & 0x80) != 0
	s.overflow = (flag & 0x40) != 0
	s.brk = (flag & 0x10) != 0
	s.irq = (flag & 0x04) != 0
	s.zero = (flag & 0x02) != 0
	s.carry = (flag & 0x01) != 0
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
