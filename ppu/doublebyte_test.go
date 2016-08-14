package ppu

import "testing"

func TestInit(t *testing.T) {
	var double doubleByte

	if double.data[0] != 0 {
		t.Error()
	}
	if double.data[1] != 0 {
		t.Error()
	}
}

func TestWrite(t *testing.T) {
	var double doubleByte
	double.Write(0xca)
	double.Write(0xfe)

	if double.data[0] != 0xca {
		t.Error()
	}

	if double.data[1] != 0xfe {
		t.Error()
	}

	if double.Value() != 0xcafe {
		t.Error()
	}
}

func TestSet(t *testing.T) {
	var double doubleByte
	double.Set(0xcafe)

	if double.data[0] != 0xca {
		t.Error()
	}

	if double.data[1] != 0xfe {
		t.Error()
	}
}
