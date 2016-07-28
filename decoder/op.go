package decoder

type Op int

const (
	LDA Op = iota
	LDX
	LDY
	STA
	STX
	STY
	TAX
	TAY
	TSX
	TXA
	TXS
	TYA
	ADC
	AND
	ASL
	BIT
	CMP
	CPX
	CPY
	DEC
	DEX
	DEY
	EOR
	INC
	INX
	INY
	LSR
	ORA
	ROL
	ROR
	SBC
	PHA
	PHP
	PLA
	PLP
	JMP
	JSR
	RTS
	RTI
	BCC
	BCS
	BEQ
	BMI
	BNE
	BPL
	BVC
	BVS
	CLC
	CLD
	CLI
	CLV
	SEC
	SED
	SEI
	BRK
	NOP
)

func opcodeName(op Op) string {
	switch op {
	case LDA:
		return "LDA"
	case LDX:
		return "LDX"
	case LDY:
		return "LDY"
	case STA:
		return "STA"
	case STX:
		return "STX"
	case STY:
		return "STY"
	case TAX:
		return "TAX"
	case TAY:
		return "TAY"
	case TSX:
		return "TSX"
	case TXA:
		return "TXA"
	case TXS:
		return "TXS"
	case TYA:
		return "TYA"
	case ADC:
		return "ADC"
	case AND:
		return "AND"
	case ASL:
		return "ASL"
	case BIT:
		return "BIT"
	case CMP:
		return "CMP"
	case CPX:
		return "CPX"
	case CPY:
		return "CPY"
	case DEC:
		return "DEC"
	case DEX:
		return "DEX"
	case DEY:
		return "DEY"
	case EOR:
		return "EOR"
	case INC:
		return "INC"
	case INX:
		return "INX"
	case INY:
		return "INY"
	case LSR:
		return "LSR"
	case ORA:
		return "ORA"
	case ROL:
		return "ROL"
	case ROR:
		return "ROR"
	case SBC:
		return "SBC"
	case PHA:
		return "PHA"
	case PHP:
		return "PHP"
	case PLA:
		return "PLA"
	case PLP:
		return "PLP"
	case JMP:
		return "JMP"
	case JSR:
		return "JSR"
	case RTS:
		return "RTS"
	case RTI:
		return "RTI"
	case BCC:
		return "BCC"
	case BCS:
		return "BCS"
	case BEQ:
		return "BEQ"
	case BMI:
		return "BMI"
	case BNE:
		return "BNE"
	case BPL:
		return "BPL"
	case BVC:
		return "BVC"
	case BVS:
		return "BVS"
	case CLC:
		return "CLC"
	case CLD:
		return "CLD"
	case CLI:
		return "CLI"
	case CLV:
		return "CLV"
	case SEC:
		return "SEC"
	case SED:
		return "SED"
	case SEI:
		return "SEI"
	case BRK:
		return "BRK"
	case NOP:
		return "NOP"
	default:
		return "UNKNOWN"
	}
}
