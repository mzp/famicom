package bits

import "testing"

func TestFlag(t *testing.T) {
	if !IsFlag(0x01, 0) {
		t.Error()
	}
	if IsFlag(0x01, 1) {
		t.Error()
	}

	if !IsFlag(0x80, 7) {
		t.Error()
	}
	if IsFlag(0x7F, 7) {
		t.Error()
	}

}
