package debug

import (
	"fmt"
	"os"

	"github.com/mzp/famicom/sprite"
)

func DumpSprite(sprites []sprite.Sprite) {
	file, err := os.Create("log/sprite.dump")
	defer file.Close()

	if err != nil {
		return
	}

	for i := 0; i < len(sprites); i++ {
		if i != 0 && i%4 == 0 {
			fmt.Fprintln(file)
		}
		sprite := sprites[i]
		fmt.Fprintf(file, "(%02x,%02x, %02x) ", sprite.X, sprite.Y, sprite.Pattern)
	}

	fmt.Fprintln(file, "\n")
}
