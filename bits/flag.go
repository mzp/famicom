package bits

func IsFlag(value byte, pos uint) bool {
	flag := byte(1) << pos
	return (value & flag) != 0
}
