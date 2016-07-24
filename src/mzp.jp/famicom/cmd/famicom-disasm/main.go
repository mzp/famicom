package main

import (
	"fmt"
	"os"

	"mzp.jp/famicom/disasm"
)

func assert(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func read(path string) {
	file, err := os.Open(path)
	assert(err)

	defer file.Close()

	for _, inst := range disasm.Disasm(file) {
		fmt.Println(inst.String())
	}
}

func main() {
	for _, path := range os.Args[1:] {
		read(path)
	}
}
