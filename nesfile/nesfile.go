package nesfile

type T struct {
	program   []byte
	character []byte
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

	return T{
		program:   get(data, headerSize, programSize),
		character: get(data, headerSize+programSize, characterSize),
	}
}
