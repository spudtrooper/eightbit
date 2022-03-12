package convert

import (
	"image/color"

	"github.com/jyotiska/go-webcolors"
)

func colorName(c color.Color) string {
	r, g, b, _ := c.RGBA()
	name := webcolors.RGBToName([]int{int(r), int(g), int(b)}, "html4")
	return name
}
