package memory

type Memory struct {
	Data      [0x10000]byte
	writeTrap map[uint16]func(byte)
	readTrap  map[uint16]func() byte
	mirror    map[uint16]uint16
}

func New() *Memory {
	var memory Memory
	memory.writeTrap = map[uint16]func(byte){}
	memory.readTrap = map[uint16]func() byte{}
	memory.mirror = map[uint16]uint16{}
	return &memory
}

func (memory *Memory) Load(address uint16, data []byte) {
	copy(memory.Data[address:], data)
}

func (memory *Memory) Read(address uint16) uint8 {
	f, ok := memory.readTrap[address]

	if ok {
		return f()
	} else {
		return memory.Data[memory.resolve(address)]
	}
}

func (memory *Memory) Read16(address uint16) uint16 {
	low := memory.Read(address)
	high := memory.Read(address + 1)
	return uint16(high)<<8 | uint16(low)
}

func (memory *Memory) ReadRange(address uint16, size uint16) []byte {
	result := make([]byte, size)

	for i := uint16(0); i < size; i++ {
		result[i] = memory.Read(address + i)
	}

	return result
}

func (memory *Memory) Write(address uint16, value byte) {
	f, ok := memory.writeTrap[address]

	if ok {
		f(value)
	} else {
		memory.Data[memory.resolve(address)] = value
	}
}

func (memory *Memory) WriteTrap(address uint16, f func(byte)) {
	memory.writeTrap[address] = f
}

func (memory *Memory) ReadTrap(address uint16, f func() byte) {
	memory.readTrap[address] = f
}

func (self *Memory) SetMirror(from, to uint16, size uint16) {
	for i := uint16(0); i < size; i++ {
		self.mirror[to+i] = from + i
	}
}

func (self *Memory) resolve(address uint16) uint16 {
	address2, ok := self.mirror[address]

	if ok {
		return address2
	} else {
		return address
	}
}
