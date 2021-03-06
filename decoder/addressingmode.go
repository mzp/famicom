package decoder

import "fmt"

type AddressingMode int

const (
	None AddressingMode = iota
	Accumlator
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Relative
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
)

func formatAddressingMode(mode AddressingMode, value int) string {
	switch mode {
	case None:
		return ""
	case Accumlator:
		return "a"
	case Immediate:
		return fmt.Sprintf("#$%04x", value)
	case ZeroPage:
		return fmt.Sprintf("<$%02x", value)
	case ZeroPageX:
		return fmt.Sprintf("<$%02x, x", value)
	case ZeroPageY:
		return fmt.Sprintf("<$%02x, y", value)
	case Relative:
		return fmt.Sprintf(" $%04x", value)
	case Absolute:
		return fmt.Sprintf(" $%04x", value)
	case AbsoluteX:
		return fmt.Sprintf(" $%04x, x", value)
	case AbsoluteY:
		return fmt.Sprintf(" $%04x, y", value)
	case Indirect:
		return fmt.Sprintf("($%04x)", value)
	case IndirectX:
		return fmt.Sprintf("($%04x, x)", value)
	case IndirectY:
		return fmt.Sprintf("($%04x), y", value)
	default:
		return ""
	}
}
