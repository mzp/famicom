package debug

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/mzp/famicom/pattern"
)

func DumpPatternImage(name string, patterns []pattern.Pattern) {
	img := image.NewRGBA(image.Rect(0, 0, pattern.WIDTH*8, pattern.HEIGHT*len(patterns)/8))

	pallets := []color.Color{
		color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
		color.RGBA{0xA0, 0xA0, 0xA0, 0xFF},
		color.RGBA{0, 0xFF, 0, 0xFF},
		color.RGBA{0, 0, 0, 0xFF},
	}

	for i, p := range patterns {
		x := (i % 8) * pattern.WIDTH
		y := (i / 8) * pattern.HEIGHT

		pattern.PutImage(img, x, y, p, pallets)
	}

	file, _ := os.Create(fmt.Sprintf("log/%s-pattern.png", name))
	defer file.Close()
	png.Encode(file, img)
}
