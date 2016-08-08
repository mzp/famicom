package pad

type Pad struct {
	memory      [8]byte
	pos         int
	resetStatus ResetStatus
}

type ResetStatus int

const (
	Init ResetStatus = iota
	Wait0
)

type Button int

const (
	A Button = iota
	B
	Select
	Start
	Up
	Down
	Left
	Right
)

func New() *Pad {
	t := Pad{}
	return &t
}

func (pad *Pad) SetButton(button Button, pressed bool) {
	if pressed {
		pad.memory[button] = 1
	} else {
		pad.memory[button] = 0
	}
}

func (pad *Pad) Read() byte {
	t := pad.memory[pad.pos]
	pad.pos = (pad.pos + 1) % 8
	return t
}

func (pad *Pad) Write(value byte) {
	pad.memory[pad.pos] = value

	switch pad.resetStatus {
	case Init:
		if value == 1 {
			pad.resetStatus = Wait0
		}
	case Wait0:
		if value == 0 {
			pad.pos = 0
			pad.resetStatus = Init
		}
	}
}
