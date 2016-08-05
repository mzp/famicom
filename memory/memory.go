package memory

type Memory struct {
	Data [0xFFFF]byte
}

func New() *Memory {
	var memory Memory
	return &memory
}

func (memory *Memory) Load(address uint16, data []byte) {
	copy(memory.Data[address:], data)
}

func (memory *Memory) Read(address uint16) uint8 {
	return memory.Data[address]
}

func (memory *Memory) Read16(address uint16) uint16 {
	low := memory.Data[address]
	high := memory.Data[address+1]
	return uint16(high)<<8 | uint16(low)
}

func (memory *Memory) ReadRange(address uint16, size uint16) []byte {
	return memory.Data[address : address+size]
}

func (memory *Memory) Write(address uint16, value byte) {
	memory.Data[address] = value
}
