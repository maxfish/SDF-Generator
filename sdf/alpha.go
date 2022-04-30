package sdf

import (
	"image"
	"image/color"
)

func AlphaFromGrey(grey image.Image) image.Image {
	rect := grey.Bounds()
	out := image.NewAlpha16(rect)
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			grey := grey.At(x, y).(color.Gray16)
			out.SetAlpha16(int(x), int(y), color.Alpha16{A: grey.Y})
		}
	}
	return out
}
