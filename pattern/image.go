package pattern

import (
	"image"
	"image/color"
)

func set(img *image.RGBA, x, y int, c color.Color) {
	min, max := img.Rect.Min, img.Rect.Max
	width := max.X - min.X
	height := max.Y - min.Y

	img.Set(x%width+min.X, y%height+min.Y, c)
}

func PutImage(img *image.RGBA, x int, y int, pattern Pattern, pallets []color.Color) {
	for py, row := range pattern {
		for px, v := range row {
			c := pallets[v]
			set(img, x+px, y+py, c)
		}
	}
}
