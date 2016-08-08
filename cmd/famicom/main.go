package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/mzp/famicom/cpu"
	"github.com/mzp/famicom/ioregister"
	"github.com/mzp/famicom/memory"
	"github.com/mzp/famicom/nesfile"
	"github.com/mzp/famicom/pad"
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
	m, c := createCPU(rom)

	pad1 := pad.New()
	pad2 := pad.New()

	ioregister.ConnectPPU(m, p)
	ioregister.ConnectPad(m, pad1, pad2)

	go run(c)

	window.CreateWindow("Famicom", func(getInput window.GetInput) image.Image {
		scanPad(pad1, pad2, getInput)
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

func createCPU(rom nesfile.T) (*memory.Memory, *cpu.CPU) {
	m := memory.New()
	m.Load(0x8000, rom.Program)
	return m, cpu.New(m, 0x8000)
}

func run(c *cpu.CPU) {
	for {
		c.Step()
	}
}

func scanPad(pad1, pad2 *pad.Pad, getInput window.GetInput) {
	pad1.SetButton(pad.A, getInput(glfw.KeyA))
	pad1.SetButton(pad.B, getInput(glfw.KeyS))
	pad1.SetButton(pad.Select, getInput(glfw.KeyZ))
	pad1.SetButton(pad.Start, getInput(glfw.KeySpace))
	pad1.SetButton(pad.Up, getInput(glfw.KeyUp))
	pad1.SetButton(pad.Down, getInput(glfw.KeyDown))
	pad1.SetButton(pad.Left, getInput(glfw.KeyLeft))
	pad1.SetButton(pad.Right, getInput(glfw.KeyRight))

	pad2.SetButton(pad.A, getInput(glfw.KeyI))
	pad2.SetButton(pad.B, getInput(glfw.KeyO))
	pad2.SetButton(pad.Select, getInput(glfw.KeyN))
	pad2.SetButton(pad.Start, getInput(glfw.KeyM))
	pad2.SetButton(pad.Up, getInput(glfw.KeyK))
	pad2.SetButton(pad.Down, getInput(glfw.KeyJ))
	pad2.SetButton(pad.Left, getInput(glfw.KeyH))
	pad2.SetButton(pad.Right, getInput(glfw.KeyL))
}
