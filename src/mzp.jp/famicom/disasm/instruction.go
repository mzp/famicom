package disasm

import "strings"

type Instruction struct {
	op int
}

func (inst Instruction) String() string {
	return strings.ToLower(opcodeName(inst.op))
}
