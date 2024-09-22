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

type DoGOptions struct {
	Sigma1  float32 // 1st Gaussian blur
	Sigma2  float32 // 2nd Gaussian blur
	Epsilon float32 // threshold
	Tau     float32 // Gaussian scalar coefficient
	// Phi     float32 // hyperbolic tangent
}

func thresHoldDoG(diff, epsilon float32) uint8 {
	if diff >= epsilon {
		return 255
	}
	// return uint8(1 + math.Tanh(float64(phi*(diff-epsilon))))
	return 0
}

// Apply Difference of Gaussians as preprocessor for edge detection
func DoG(img image.Image, opts DoGOptions) image.Image {
	b1 := gift.New(gift.GaussianBlur(opts.Sigma1))
	b2 := gift.New(gift.GaussianBlur(opts.Sigma2))

	dst1 := image.NewGray(b1.Bounds(img.Bounds()))
	b1.Draw(dst1, img)

	dst2 := image.NewGray(b2.Bounds(img.Bounds()))
	b2.Draw(dst2, img)

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	doG := image.NewGray(img.Bounds())
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			p1 := dst1.GrayAt(j, i)
			p2 := dst2.GrayAt(j, i)

			diff := thresHoldDoG((1+opts.Tau)*float32(p1.Y)-opts.Tau*float32(p2.Y), opts.Epsilon)
			doG.Set(j, i, color.Gray{Y: diff})
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
