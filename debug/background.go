package debug

import (
	"fmt"
	"os"
)

func DumpBackground(n int, nameTable []byte, attributeTable []byte) {
	file, err := os.Create(fmt.Sprintf("log/bg%d.dump", n))
	defer file.Close()

	if err != nil {
		return
	}

	for i := 0; i < len(nameTable); i++ {
		if i != 0 && i%32 == 0 {
			fmt.Fprintln(file)
		}

		fmt.Fprintf(file, "%02x ", nameTable[i])
	}
	fmt.Fprintln(file, "\n")

	for i := 0; i < len(attributeTable); i++ {
		if i != 0 && i%32 == 0 {
			fmt.Fprintln(file)
		}

		fmt.Fprintf(file, "% 2x ", attributeTable[i])
	}
	fmt.Fprintln(file, "\n")
}
