package nesfile

import "github.com/mzp/famicom/bits"

type T struct {
	Program        []byte
	Character      []byte
	VerticalMirror bool
}

const (
	headerSize        = 16
	programSizeUnit   = 16 * 1024
	characterSizeUnit = 8 * 1024
)

func get(data []byte, start, length int) []byte {
	return data[start : start+length]
}

func Load(data []byte) T {
	programSize := int(data[4]) * programSizeUnit
	characterSize := int(data[5]) * characterSizeUnit
	flag := data[6]

	return T{
		Program:        get(data, headerSize, programSize),
		Character:      get(data, headerSize+programSize, characterSize),
		VerticalMirror: bits.IsFlag(flag, 0),
	}
}
