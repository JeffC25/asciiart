package asciiart

import (
	"image"
	"image/color"

	"github.com/disintegration/gift"
)

type Edge int

const (
	None         Edge = iota
	Horizontal        // "_"
	Vertical          // "|"
	DiagonalUp        // "/"
	DiagonalDown      // "\"
)

func DoG(img image.Image, sigma1, sigma2 float32) image.Image {
	b1 := gift.New(gift.GaussianBlur(sigma1))
	b2 := gift.New(gift.GaussianBlur(sigma2))

	dst1 := image.NewGray(img.Bounds())
	b1.Draw(dst1, img)

	dst2 := image.NewGray(img.Bounds())
	b2.Draw(dst2, img)

	doG := image.NewGray(img.Bounds())
	for i := 0; i < img.Bounds().Dy(); i++ {
		for j := 0; j < img.Bounds().Dx(); j++ {
			p1 := dst1.GrayAt(j, i)
			p2 := dst2.GrayAt(j, i)

			diff := max(p1.Y-p2.Y, p2.Y-p1.Y)
			doG.SetGray(j, i, color.Gray{Y: diff})
		}

	}
	return doG
}
