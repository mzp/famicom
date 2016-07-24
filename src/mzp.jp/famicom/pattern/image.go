package pattern

import (
	"image"
	"image/color"
)

func PutImage(img *image.RGBA, x int, y int, pattern Pattern, pallets []color.Color) {
	for py, row := range pattern {
		for px, v := range row {
			c := pallets[v]
			img.Set(x+px, y+py, c)
		}
	}
}
