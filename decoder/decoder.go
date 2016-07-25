package decoder

func readValue(data []byte, pc, size int) int {
	value := 0

	for i := 0; i < size; i++ {
		v := data[pc+i]
		value |= int(v) << uint(8*i)
	}

	return value
}

func Decode(data []byte, pc int) (Instruction, int) {
	op := data[pc]
	entry := decodeTable[op]

	value := readValue(data, pc+1, entry.size)

	inst := Instruction{entry.op, entry.addressingMode, value}
	return inst, entry.size + 1
}

func DecodeAll(data []byte) []Instruction {
	instructions := []Instruction{}

	pc := 0

	for pc < len(data) {
		inst, n := Decode(data, pc)
		pc += n
		instructions = append(instructions, inst)
	}

	return instructions
}
