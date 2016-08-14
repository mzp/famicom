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

	if len(file.Program) != 2*16*1024 {
		t.Errorf("unmatch program size: %x", len(file.Program))
	}

	if len(file.Character) != 8*1024 {
		t.Errorf("unmatch character size: %x", len(file.Character))
	}

	if !startWith(file.Program, []byte{0xca, 0xfe, 0xba, 0xbe}) {
		t.Errorf("cannot load program rom: %v", file.Program[0:16])
	}

	if !startWith(file.Character, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Errorf("cannot load character rom: %v", file.Character[0:16])
	}

	if !file.VerticalMirror {
		t.Error("cannot read vertical mirror flag")
	}
}
