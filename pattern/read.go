package pattern

import "io"

func bits(x byte) []byte {
	xs := make([]byte, 8)

	for i := 0; i < len(xs); i++ {
		xs[i] = (x & 0x80) >> 7
		x = x << 1
	}

	return xs
}

func pattern(low, high []byte) Pattern {
	result := Pattern{}

	for y := 0; y < HEIGHT; y++ {
		l, h := bits(low[y]), bits(high[y])

		for x := 0; x < WIDTH; x++ {
			result[y][x] = l[x] + h[x]<<1
		}
	}

	return result
}

func ReadFromBytes(data []byte) (Pattern, bool) {
	low, high := data[0:UNIT_BYTES], data[UNIT_BYTES:UNIT_BYTES*2]

	if len(low) == UNIT_BYTES && len(high) == UNIT_BYTES {
		return pattern(low, high), true
	} else {
		return Pattern{}, false
	}
}

func ReadAll(reader io.Reader) []Pattern {
	patterns := []Pattern{}

	buf := make([]byte, READ_BYTES)

	for {
		n, _ := reader.Read(buf)
		if n != READ_BYTES {
			break
		}

		x, ok := ReadFromBytes(buf)

		if !ok {
			break
		}

		patterns = append(patterns, x)
	}

	return patterns
}

func ReadAllFromBytes(data []byte) []Pattern {
	patterns := []Pattern{}

	for i := 0; i < len(data); i += READ_BYTES {
		x, ok := ReadFromBytes(data[i:])

		if !ok {
			break
		}

		patterns = append(patterns, x)
	}
	return patterns
}
