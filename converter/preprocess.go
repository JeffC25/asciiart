package asciiart

import (
	"image"

	"github.com/disintegration/gift"
)

// Grayscale and Downscale
func Preprocess(img image.Image, factor int) image.Image {
	g := gift.New(gift.Resize(img.Bounds().Dx()/factor, img.Bounds().Dy()/factor, gift.BoxResampling))
	dst := image.NewGray(g.Bounds(img.Bounds()))
	g.Draw(dst, img)
	return dst
}
