package memory

type Memory struct {
  Data [0xFFFF]byte
}

type Address int

func New() *Memory {
	var memory Memory
	return &memory
}

func (memory *Memory) Load(address Address, data []byte) {
	copy(memory.Data[address:], data)
}

func (memory *Memory) Get(address Address) byte {
	return memory.Data[address]
}
