package disasm

import "io"

func readByte(reader io.Reader) (byte, error) {
	data := make([]byte, 1)
	_, err := reader.Read(data)
	return data[0], err
}

func readValue(reader io.Reader, n int) (int, error) {
	value := 0
	for i := 0; i < n; i++ {
		v, err := readByte(reader)

		if err != nil {
			return 0, err
		}

		value += int(v) << uint(8*i)
	}

	return value, nil
}

func Disasm(reader io.Reader) []Instruction {
	instructions := []Instruction{}

	for {
		op, err := readByte(reader)

		if err != nil {
			break
		}

		entry := decodeTable[op]

		value, err := readValue(reader, entry.size)

		if err != nil {
			break
		}

		instructions = append(
			instructions,
			Instruction{entry.op, entry.addressingMode, value})
	}

	return instructions
}
