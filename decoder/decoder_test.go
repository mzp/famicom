package decoder

import "testing"

func TestDecode(t *testing.T) {

	data := []byte{
	  0xea, // EOP
	  0xad, // lda $2000
	  0x00,
	  0x20,
	}
	inst, size := Decode(data, 1)

	if size != 3 {
		t.Error()
	}

	if inst.Op != LDA {
		t.Error()
	}

	if inst.AddressingMode != Absolute {
		t.Error()
	}

	if inst.Value != 0x2000 {
		t.Error()
	}
}
