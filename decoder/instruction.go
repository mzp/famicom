package decoder

import (
	"fmt"
	"strings"
)

type Instruction struct {
	op             op
	addressingMode addressingMode
	value          int
}

func (inst Instruction) String() string {
	if inst.addressingMode == None {
		return strings.ToLower(opcodeName(inst.op))
	} else {
		return fmt.Sprintf(
			"%s %s",
			strings.ToLower(opcodeName(inst.op)),
			formatAddressingMode(inst.addressingMode, inst.value))
	}
}
