package sprite

import "testing"

func TestSpriteSize(t *testing.T) {
	sm := New()

	sprites := sm.Get()

	if len(sprites) != 64 {
		t.Errorf("cannot get 64 sprites: %d", len(sprites))
	}
}

func TestSprite(t *testing.T) {
	sm := New()
	sm.SetAddress(0)
	sm.Write(40)
	sm.Write(1)
	sm.Write(0xe1)
	sm.Write(30)

	expect := Sprite{
		X:               30,
		Y:               40,
		Pattern:         1,
		Palette:         1,
		FlipVertical:    true,
		FlipHorizon:     true,
		UnderBackground: true,
	}

	sprites := sm.Get()
	if sprites[0] != expect {
		t.Errorf("Sprite %v", sprites[0])
	}
}

func TestMultiSprite(t *testing.T) {
	sm := New()

	sm.SetAddress(0)
	sm.Write(40)
	sm.Write(1)
	sm.Write(0xe1)
	sm.Write(30)

	sm.Write(40)
	sm.Write(1)
	sm.Write(0xe1)
	sm.Write(30)

	expect := Sprite{
		X:               30,
		Y:               40,
		Pattern:         1,
		Palette:         1,
		FlipVertical:    true,
		FlipHorizon:     true,
		UnderBackground: true,
	}

	sprites := sm.Get()
	if sprites[0] != expect {
		t.Errorf("Sprite %v", sprites[0])
	}
	if sprites[1] != expect {
		t.Errorf("Sprite %v", sprites[1])
	}
}

func TestCopy(t *testing.T) {
	sm := New()

	var data = [255]byte{
		40, 1, 0xe1, 30,
		40, 1, 0xe1, 30,
	}

	sm.Copy(data[:])

	expect := Sprite{
		X:               30,
		Y:               40,
		Pattern:         1,
		Palette:         1,
		FlipVertical:    true,
		FlipHorizon:     true,
		UnderBackground: true,
	}

	sprites := sm.Get()
	if sprites[0] != expect {
		t.Errorf("Sprite %v", sprites[0])
	}
	if sprites[1] != expect {
		t.Errorf("Sprite %v", sprites[1])
	}
}
