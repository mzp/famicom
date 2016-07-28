package decoder

import (
	"fmt"
	"strings"
)

type Instruction struct {
	Op             Op
	AddressingMode AddressingMode
	Value          int
}

func (inst Instruction) String() string {
	if inst.AddressingMode == None {
		return strings.ToLower(opcodeName(inst.Op))
	} else {
		return fmt.Sprintf(
			"%s %s",
			strings.ToLower(opcodeName(inst.Op)),
			formatAddressingMode(inst.AddressingMode, inst.Value))
	}
}
