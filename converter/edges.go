package asciiart

import (
	"image"
	"image/color"
	"math"

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

// Apply Difference of Gaussians as preprocessor for edge detection
func DoG(img image.Image, sigma1, sigma2 float32) image.Image {
	b1 := gift.New(gift.GaussianBlur(sigma1))
	b2 := gift.New(gift.GaussianBlur(sigma2))

	dst1 := image.NewRGBA(img.Bounds())
	b1.Draw(dst1, img)

	dst2 := image.NewRGBA(img.Bounds())
	b2.Draw(dst2, img)

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	doG := image.NewGray(img.Bounds())
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			p1 := dst1.RGBAAt(j, i)
			p2 := dst2.RGBAAt(j, i)

			diffR := uint8(math.Abs(float64(p1.R) - float64(p2.R)))
			diffG := uint8(math.Abs(float64(p1.G) - float64(p2.G)))
			diffB := uint8(math.Abs(float64(p1.B) - float64(p2.B)))

			doG.Set(j, i, color.RGBA{R: diffR, G: diffG, B: diffB})
		}

	}
	return doG
}

// Map angle values to discrete edges
func mapAngleToEdge(angle, Gx, Gy float64) Edge {
	if Gx == 0 && Gy == 0 {
		return None
	}

	// Convert angle to degrees for easier mapping
	degrees := angle * (180.0 / math.Pi)

	// Map angle to edge type
	switch {
	case degrees > -22.5 && degrees <= 22.5:
		return Horizontal // Close to 0 degrees, horizontal edge
	case degrees > 22.5 && degrees <= 67.5:
		return DiagonalDown // Between 22.5 and 67.5 degrees, diagonal down
	case degrees > 67.5 || degrees <= -67.5:
		return Vertical // Close to 90 or -90 degrees, vertical edge
	case degrees > -67.5 && degrees <= -22.5:
		return DiagonalUp // Between -22.5 and -67.5 degrees, diagonal up
	default:
		return None
	}
}

func Sobel(img image.Image) image.Image {

	g := gift.New(gift.Sobel())
	dst := image.NewGray(g.Bounds(img.Bounds()))
	g.Draw(dst, img)
	return dst
}

func SobelX(img image.Image) image.Image {
	// Create the horizontal gradient filter (Gx)
	sobelGx := gift.New(
		gift.Convolution([]float32{
			-1, 0, 1,
			-2, 0, 2,
			-1, 0, 1,
		}, false, false, false, 0.0),
	)

	// Prepare destination images
	dstGx := image.NewGray(sobelGx.Bounds(img.Bounds()))

	// Apply the filter
	sobelGx.Draw(dstGx, img)
	return dstGx
}
func SobelY(img image.Image) image.Image {
	// Create the vertical gradient filter (Gy)
	sobelGy := gift.New(
		gift.Convolution([]float32{
			-1, -2, -1,
			0, 0, 0,
			1, 2, 1,
		}, false, false, false, 0.0),
	)

	// Prepare destination images
	dstGy := image.NewGray(sobelGy.Bounds(img.Bounds()))

	// Apply the filter
	sobelGy.Draw(dstGy, img)
	return dstGy
}
