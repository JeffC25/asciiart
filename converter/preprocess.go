package asciiart

import (
	"image"
	"math"

	"github.com/disintegration/gift"
)

// Grayscale and Downscale
func Preprocess(img image.Image, width int) *image.Gray {
	scale := float64(img.Bounds().Dx()) / float64(width)
	height := int(math.Ceil(float64(img.Bounds().Dy()) / scale))

	g := gift.New(gift.Resize(width, height, gift.BoxResampling))
	dst := image.NewGray(g.Bounds(img.Bounds()))
	g.Draw(dst, img)

	return dst
}
