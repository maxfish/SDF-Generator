package sdf

import (
	"math"

	"image"
	"image/color"
)

func GenerateDistanceFieldImage(inputImage image.Image, downscale int, spread float64, sourceChannels [4]bool, threshold float64, crop bool) image.Image {
	grid := gridFromImage(inputImage, sourceChannels, threshold, crop)
	return sdfFromGrid(grid, downscale, spread)
}

func gridFromImage(inputImage image.Image, sourceChannels [4]bool, threshold float64, crop bool) *BoolGrid {
	inputWidth, inputHeight := inputImage.Bounds().Max.X, inputImage.Bounds().Max.Y
	grid := NewBoolGrid(inputWidth, inputHeight)
	intThreshold := uint32(threshold * math.MaxUint16)

	for y := 0; y < inputHeight; y++ {
		for x := 0; x < inputWidth; x++ {
			red, green, blue, alpha := inputImage.At(x, y).RGBA()
			// The image contains "something" at x,y if one of the selected channels has a value
			// greater than the specified threshold
			if (sourceChannels[0] && red >= intThreshold) ||
				(sourceChannels[1] && green >= intThreshold) ||
				(sourceChannels[2] && blue >= intThreshold) ||
				(sourceChannels[3] && alpha >= intThreshold) {
				grid.Set(true, x, y)
			}
		}
	}

	if crop {
		return grid.Crop()
	}

	return grid
}

func sdfFromGrid(grid *BoolGrid, downscale int, spread float64) image.Image {
	gridWidth, gridHeight := grid.W, grid.H
	delta := int(math.Floor(spread))

	// The output image is slightly bigger since it needs space for the spread
	outputWidth := gridWidth/downscale + delta*2
	outputHeight := gridHeight/downscale + delta*2
	outputImage := image.NewGray16(image.Rect(0, 0, int(outputWidth), int(outputHeight)))

	for y := 0; y < outputHeight; y++ {
		for x := 0; x < outputWidth; x++ {
			centerX := (x-delta)*downscale + downscale/2
			centerY := (y-delta)*downscale + downscale/2

			signedDistance := findSignedDistance(grid, centerX, centerY, spread)

			// Convert the distance into a pixel value
			value := 0.5 + 0.5*(signedDistance/spread)
			value = math.Min(1, math.Max(0.0, value)) * math.MaxUint16
			outputImage.SetGray16(int(x), int(y), color.Gray16{Y: uint16(value)})
		}
	}

	return outputImage
}

func findSignedDistance(grid *BoolGrid, centerX, centerY int, spread float64) float64 {
	gridValue := grid.At(centerX, centerY)
	delta := int(math.Floor(spread))
	closestDistance := delta * delta

	for y := -delta; y <= delta; y++ {
		for x := -delta; x <= delta; x++ {
			pointX := centerX + x
			pointY := centerY + y
			if gridValue != grid.At(pointX, pointY) {
				distance := squaredDistance(centerX, centerY, pointX, pointY)
				if distance < closestDistance {
					closestDistance = distance
				}
			}
		}
	}

	closestDist := math.Min(math.Sqrt(float64(closestDistance)), spread)
	if gridValue {
		// Inside
		return 1.0 * closestDist
	} else {
		// Outside
		return -1.0 * closestDist
	}
}

func squaredDistance(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}
