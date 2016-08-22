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

func PutImageWithFlip(img *image.RGBA, x int, y int, pattern Pattern, pallets []color.Color, flipH, flipV bool, skipZero bool) {
	for py := 0; py < HEIGHT; py++ {
		for px := 0; px < WIDTH; px++ {
			var row [8]byte
			var v byte

			if flipV {
				row = pattern[HEIGHT-py-1]
			} else {
				row = pattern[py]
			}

			if flipH {
				v = row[WIDTH-px-1]
			} else {
				v = row[px]
			}

			if v != 0 || !skipZero {
				c := pallets[v]
				set(img, x+px, y+py, c)
			}
		}
	}
}

func PutImage(img *image.RGBA, x int, y int, pattern Pattern, pallets []color.Color) {
	PutImageWithFlip(img, x, y, pattern, pallets, false, false, false)
}
