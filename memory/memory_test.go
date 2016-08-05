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
