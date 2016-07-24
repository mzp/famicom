package disasm

import "io"

func readByte(reader io.Reader) (byte, error) {
	data := make([]byte, 1)

	_, err := reader.Read(data)
	return data[0], err
}

func Disasm(reader io.Reader) []Instruction {
	instructions := []Instruction{}

	for {
		op, err := readByte(reader)

		if err != nil {
			break
		}

		entry := decodeTable[op]

		for i := 0; i < entry.size; i++ {
			readByte(reader)
		}

		instructions = append(
			instructions,
			Instruction{entry.op})
	}

	return instructions
}
