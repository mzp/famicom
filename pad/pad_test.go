package pad

import (
	"testing"

	"github.com/mzp/famicom/bits"
)

func TestRead1(t *testing.T) {
	pad := New()
	pad.SetButton(A, true)
	pad.SetButton(Select, true)

	if !bits.IsFlag(pad.Read(), 0) {
		t.Error("cannot read a button")
	}
	if bits.IsFlag(pad.Read(), 0) {
		t.Error("cannot read b button")
	}
	if !bits.IsFlag(pad.Read(), 0) {
		t.Error("cannot read select button")
	}
}

func TestReadLoop(t *testing.T) {
	pad := New()
	pad.SetButton(A, true)

	for i := 0; i < 8; i++ {
		pad.Read()
	}
	if !bits.IsFlag(pad.Read(), 0) {
		t.Error("cannot read a button")
	}
}

func TestWrite(t *testing.T) {
	pad := New()

	pad.SetButton(A, true)

	pad.Write(0)

	if bits.IsFlag(pad.Read(), 0) {
		t.Error("cannot read a button")
	}

	pad.Write(0)
	pad.SetButton(B, true)
	if !bits.IsFlag(pad.Read(), 0) {
		t.Error("cannot read b button", pad)
	}
}
