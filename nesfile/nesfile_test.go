package nesfile

import "testing"

func startWith(xs, ys []byte) bool {
	if len(xs) < len(ys) {
		return false
	}

	for n, y := range ys {
		x := xs[n]
		if x != y {
			return false
		}
	}
	return true
}

func TestLoad(t *testing.T) {
	var data [16 + 2*16*1024 + 1*8*1024]byte

	copy(data[0:16], []byte{
		0x4e, 0x45, 0x53, 0x1a, // NES
		0x02, // PRG Size
		0x01, // Chr Size
		0x01, // Vertical Mirror
		0x00, // Mapper 0
	})
	copy(data[16:], []byte{
		0xca, 0xfe, 0xba, 0xbe,
	})
	copy(data[16+2*16*1024:], []byte{
		0xde, 0xad, 0xbe, 0xef,
	})

	file := Load(data[:])

	if len(file.program) != 2*16*1024 {
		t.Errorf("unmatch program size: %x", len(file.program))
	}

	if len(file.character) != 8*1024 {
		t.Errorf("unmatch character size: %x", len(file.character))
	}

	if !startWith(file.program, []byte{0xca, 0xfe, 0xba, 0xbe}) {
		t.Errorf("cannot load program rom: %v", file.program[0:16])
	}

	if !startWith(file.character, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Errorf("cannot load character rom: %v", file.character[0:16])
	}
}
