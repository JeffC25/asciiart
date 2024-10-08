package asciiart

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/disintegration/gift"
)

// Grayscale and Downscale
func GrayDownscale(img image.Image, width int, squash float32) *image.Gray {
	scale := float64(img.Bounds().Dx()) / float64(width) * float64(squash)
	height := int(math.Floor(float64(img.Bounds().Dy()) / scale))

	g := gift.New(gift.Resize(width, height, gift.BoxResampling))
	dst := image.NewGray(g.Bounds(img.Bounds()))
	g.Draw(dst, img)

	return dst
}

// Converts a grayscale image to ASCII art
func ConvertToASCIIArt(img image.Image, charset []rune) [][]rune {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	asciiArt := make([][]rune, height)
	for y := 0; y < height; y++ {
		row := make([]rune, width)
		for x := 0; x < width; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			luminance := grayColor.Y
			index := int(luminance) * (len(charset) - 1) / 255
			row[x] = charset[index]
		}
		asciiArt[y] = row
	}

	return asciiArt
}

// Prints the 2D ASCII art to the console
func PrintASCIIArt(asciiArt [][]rune) {
	for _, row := range asciiArt {
		for _, char := range row {
			fmt.Print(string(char))
		}
		fmt.Println()
	}
}
