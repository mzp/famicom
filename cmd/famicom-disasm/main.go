package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mzp/famicom/decoder"
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

	data, err := ioutil.ReadAll(file)
	assert(err)

	for _, inst := range decoder.DecodeAll(data) {
		fmt.Println(inst.String())
	}
}

func main() {
	for _, path := range os.Args[1:] {
		read(path)
	}
}
