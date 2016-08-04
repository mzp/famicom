package main

import (
	"image"
	_ "image/png"
	"os"
	"runtime"

	"github.com/mzp/famicom/window"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	img, err := load("square.png")

	if err != nil {
		panic(err)
	}

	window.CreateWindow("Famicom", func() image.Image {
		return img
	})
}

func load(file string) (image.Image, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}
