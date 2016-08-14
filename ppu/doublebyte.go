package ppu

type doubleByte struct {
	data  [2]byte
	index int
}

func (self *doubleByte) Write(value byte) {
	self.data[self.index] = value
	self.index = (self.index + 1) % 2
}

func (self *doubleByte) Value() uint16 {
	high := uint16(self.data[0])
	low := uint16(self.data[1])
	return high<<8 | low
}

func (self *doubleByte) Set(value uint16) {
	self.data[0] = uint8(value >> 8)
	self.data[1] = uint8(value)
}

func (self *doubleByte) Reset() {
	self.index = 0
}
