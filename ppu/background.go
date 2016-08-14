package ppu

import (
	"image"

	"github.com/mzp/famicom/pattern"
)

func getAttribute(attributeTable []byte, x, y int) byte {
	attribute := attributeTable[x/4+y/4*8]

	x_, y_ := x%4, y%4
	index := (x_ / 2) + (y_/2)*2

	return (attribute >> uint(index*2)) & 0x3
}

func renderBackground(ppu *PPU, nameTables, attributeTables [4][]byte) *image.RGBA {
	background := image.NewRGBA(image.Rect(0, 0, WIDTH*2, HEIGHT*2))

	for i := 0; i < 4; i++ {
		nameTable := nameTables[i]
		attributeTable := attributeTables[i]

		bx := (i % 2) * WIDTH
		by := (i / 2) * HEIGHT

		for n, v := range nameTable {
			x := n % 32
			y := n / 32

			paletteIndex := getAttribute(attributeTable, x, y)
			pattern.PutImage(background,
				bx+x*8, by+y*8,
				ppu.patterns[ppu.backgroundIndex][v],
				ppu.bgPalettes[paletteIndex])
		}
	}

	return background
}
