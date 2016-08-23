package pattern

import (
	"image"
	"image/color"
)

const (
	// size for pixel
	WIDTH  = 8
	HEIGHT = 8

	// size for byte
	UNIT_BYTES = 8
	READ_BYTES = UNIT_BYTES * 2
)

type Pattern [WIDTH][HEIGHT]byte

func set(img *image.RGBA, x, y int, c color.Color) {
	min, max := img.Rect.Min, img.Rect.Max
	width := max.X - min.X
	height := max.Y - min.Y

	img.Set(x%width+min.X, y%height+min.Y, c)
}

type Option struct {
	FlipH, FlipV  bool
	BackdropColor color.Color
}

func (self *Pattern) Put(img *image.RGBA, pos image.Point, palettes []color.Color, option Option) {
	for py := 0; py < HEIGHT; py++ {
		for px := 0; px < WIDTH; px++ {
			var row [8]byte
			var v byte

			if option.FlipV {
				row = self[HEIGHT-py-1]
			} else {
				row = self[py]
			}

			if option.FlipH {
				v = row[WIDTH-px-1]
			} else {
				v = row[px]
			}

			if v == 0 {
				if option.BackdropColor != nil {
					set(img, pos.X+px, pos.Y+py, option.BackdropColor)
				}
			} else {
				c := palettes[v]
				set(img, pos.X+px, pos.Y+py, c)
			}
		}
	}
}
