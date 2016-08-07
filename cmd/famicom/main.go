package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/mzp/famicom/cpu"
	"github.com/mzp/famicom/ioregister"
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
	rom := load(os.Args[1])
	p := createPPU(rom)
	c := createCPU(rom, p)

	go run(c)

	window.CreateWindow("Famicom", func() image.Image {
		return p.Render()
	})
}

func load(path string) nesfile.T {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	return nesfile.Load(data)
}

func createPPU(rom nesfile.T) *ppu.PPU {
	m := memory.New()
	m.Load(0x0, rom.Character)
	return ppu.New(m)
}

func createCPU(rom nesfile.T, ppu *ppu.PPU) *cpu.CPU {
	m := memory.New()
	ioregister.Connect(m, ppu)
	m.Load(0x8000, rom.Program)
	return cpu.New(m, 0x8000)
}

func run(c *cpu.CPU) {
	for {
		c.Step()
	}
}
