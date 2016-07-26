package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mzp/famicom/cpu"
	"github.com/mzp/famicom/memory"
)

func assert(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func read(path string) []byte {
	file, err := os.Open(path)
	assert(err)

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	assert(err)

	return data
}

func main() {
	data := read(os.Args[1])

	m := memory.New()
	m.Load(0x8000, data)

	c := cpu.New(m, 0x8000)

	for {
		fmt.Printf("next: %s\n%s\n", c.CurrentInstruction().String(), c.String())

		var s string

		fmt.Print("> ")
		fmt.Scan(&s)

		switch s {
		case "step":
			c.Step()
		case "q":
			os.Exit(0)
		}

		fmt.Println("")
	}
}
