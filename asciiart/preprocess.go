package asciiart

import (
	"image"
	"math"

	"github.com/disintegration/gift"
)

// Downscale image horizontally to compensate for ascii output ratios
func Squash(img image.Image, scale float32) image.Image {
	width := img.Bounds().Dx()
	height := int(math.Round(float64(img.Bounds().Dy()) / float64(scale)))

	g := gift.New(gift.Resize(width, height, gift.BoxResampling))
	dst := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dst, img)

	return dst
}
