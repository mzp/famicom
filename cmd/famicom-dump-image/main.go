package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/mzp/famicom/pattern"
)

func assert(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func read(path string) []pattern.Pattern {
	file, err := os.Open(path)
	assert(err)

	defer file.Close()

	return pattern.ReadAll(file)
}

func createImage(name string, patterns []pattern.Pattern) {
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

		p.Put(img, image.Point{x, y}, pallets, pattern.Option{})
	}

	file, err := os.Create(name + ".png")
	assert(err)
	png.Encode(file, img)
}

func name(path string) string {
	ext := filepath.Ext(path)
	basename := filepath.Base(path)

	return basename[0 : len(basename)-len(ext)]
}

func main() {
	for _, path := range os.Args[1:] {
		patterns := read(path)
		createImage(name(path), patterns)
	}
}
