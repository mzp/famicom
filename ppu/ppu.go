package ppu

import (
	"image"

	memlib "github.com/mzp/famicom/memory"
)

type PPU struct {
	memory *memlib.Memory
}

func New(m *memlib.Memory) *PPU {
	t := PPU{memory: m}
	return &t
}

func (*PPU) Render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 240))
	return img
}
