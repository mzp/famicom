package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/nesfile"
	"github.com/mzp/famicom/ppu"
	"github.com/mzp/famicom/window"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	ppu := create(os.Args[1], os.Args[2])
	ppu.SetControl2(0x08)
	window.CreateWindow("Famicom", func() image.Image {
		return ppu.Render()
	})
}

func create(path, title string) *ppu.PPU {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	rom := nesfile.Load(data)

	m := memory.New()
	m.Load(0x0, rom.Character)

	for n, c := range title {
		m.Write(0x2021+uint16(n), byte(c))
	}

	return ppu.New(m)
}
