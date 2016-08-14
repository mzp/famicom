package debug

import (
	"fmt"
	"os"
)

func DumpNameTale(n int, data []byte) {
	file, err := os.Create(fmt.Sprintf("log/nametable%d.dump", n))
	defer file.Close()

	if err != nil {
		return
	}

	for i := 0; i < len(data); i++ {
		if i != 0 && i%32 == 0 {
			fmt.Fprintln(file)
		}

		fmt.Fprintf(file, "%02x ", data[i])
	}
	fmt.Fprintln(file)
}
