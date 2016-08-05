package pattern

import "testing"

func TestRead(t *testing.T) {
	data := []byte{
		// low
		0xC0,
		0xC0,
		0xC0,
		0xF0,
		0xF8,
		0xDC,
		0xCE,
		0x00,

		// high
		0x06,
		0x0c,
		0x18,
		0x70,
		0x38,
		0x1c,
		0x0e,
		0x00,
	}
	pattern, ok := ReadFromBytes(data)

	if !ok {
		t.Error()
	}

	if pattern != [8][8]byte{
		[8]byte{1, 1, 0, 0, 0, 2, 2, 0},
		[8]byte{1, 1, 0, 0, 2, 2, 0, 0},
		[8]byte{1, 1, 0, 2, 2, 0, 0, 0},
		[8]byte{1, 3, 3, 3, 0, 0, 0, 0},
		[8]byte{1, 1, 3, 3, 3, 0, 0, 0},
		[8]byte{1, 1, 0, 3, 3, 3, 0, 0},
		[8]byte{1, 1, 0, 0, 3, 3, 3, 0},
		[8]byte{0, 0, 0, 0, 0, 0, 0, 0},
	} {
		t.Error()
	}
}

func TestReadAllFromBytes(t *testing.T) {
	data := []byte{
		// Pattern-1
		0xC0,
		0xC0,
		0xC0,
		0xF0,
		0xF8,
		0xDC,
		0xCE,
		0x00,
		0x06,
		0x0c,
		0x18,
		0x70,
		0x38,
		0x1c,
		0x0e,
		0x00,

		// Pattern-2
		0xC0,
		0xC0,
		0xC0,
		0xF0,
		0xF8,
		0xDC,
		0xCE,
		0x00,
		0x06,
		0x0c,
		0x18,
		0x70,
		0x38,
		0x1c,
		0x0e,
		0x00,
	}
	patterns := ReadAllFromBytes(data)

	if len(patterns) != 2 {
		t.Errorf("read %d pattern", len(patterns))
	}
}
