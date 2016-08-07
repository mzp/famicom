package memory

import "testing"

func TestLoad(t *testing.T) {
	memory := New()
	memory.Load(0x2000, []byte{
		0xca,
		0xfe,
	})

	if memory.Read(0x2000) != 0xca {
		t.Error()
	}
	if memory.Read(0x2001) != 0xfe {
		t.Error()
	}
}

func TestReadWrite(t *testing.T) {
	memory := New()
	memory.Write(0x2000, 42)

	if memory.Read(0x2000) != 42 {
		t.Error()
	}
}

func TestRead16(t *testing.T) {
	memory := New()
	memory.Write(0x2000, 0xfe)
	memory.Write(0x2001, 0xca)

	if memory.Read16(0x2000) != 0xcafe {
		t.Error()
	}
}

func TestReadRange(t *testing.T) {
	memory := New()
	memory.Write(0x2000, 0xfe)
	memory.Write(0x2001, 0xca)

	data := memory.ReadRange(0x2000, 2)

	if len(data) != 2 {
		t.Error()
	}

	if data[0] != 0xfe && data[1] != 0xca {
		t.Error()
	}
}

func TestWriteTrap(t *testing.T) {
	memory := New()
	n := 0

	memory.WriteTrap(0x2000, func(value byte) {
		n += int(value)
	})

	memory.Write(0x2000, 1)

	if n != 1 {
		t.Error("cannot invoke write trap")
	}

	memory.Write(0x2001, 0xff)
	if n != 1 {
		t.Error("cannot over-invoke write trap")
	}
}
