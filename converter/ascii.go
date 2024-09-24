package asciiart

import (
	"fmt"
	"image"
	"image/color"
)

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
