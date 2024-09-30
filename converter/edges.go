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
	Phi     float32 // hyperbolic tangent
}

// Extended thresholding function for DoG output
func thresholdDoG(diff, epsilon, phi float32) uint8 {
	u := diff / 255
	if u >= epsilon {
		return 255
	}
	return uint8((1 + math.Tanh(float64(phi*(u-epsilon)))) * 255)
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

			// Winnem¨oller's XDoG operator: (1 + τ) * G - τ * G
			diff := thresholdDoG((1+opts.Tau)*float32(p1.Y)-opts.Tau*float32(p2.Y), opts.Epsilon, opts.Phi)
			doG.Set(j, i, color.Gray{Y: diff})
		}
	}
	return doG
}

// Compute angle of X Y gradients and map to discrete edges if magnitude above threshold
func XYToEdge(x, y, threshold float64) Edge {
	magnitude := math.Hypot(x, y)
	if magnitude < threshold {
		return None
	}

	// math.Atan2 outputs -Pi radians to Pi radians
	angle := math.Atan2(y, x)

	// Normalize angle to the range [0, Pi]
	angle = math.Mod(angle, 2*math.Pi)
	if angle < 0 {
		angle += math.Pi
	}

	// Map the angle to the appropriate edge type
	switch {
	case angle >= 0 && angle < math.Pi/8 || angle >= 7*math.Pi/8:
		return Horizontal
	case angle >= math.Pi/8 && angle < 3*math.Pi/8:
		return DiagonalUp
	case angle >= math.Pi/2 && angle < 5*math.Pi/8:
		return Vertical
	default:
		return DiagonalDown
	}
}

// Map an image to a 2d array of Edge types
func MapEdges(img *image.Gray, sobelThreshold float64) [][]Edge {
	threshold := sobelThreshold * math.Hypot(255*4, 255*4)

	Gx := [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	Gy := [3][3]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	edges := make([][]Edge, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]Edge, width)
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// High horizontal change = vertical edge
			// High vertical change = horizontal edge
			sumX := 0
			sumY := 0
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					pixel := int(img.GrayAt(x+kx, y+ky).Y)
					sumX += pixel * Gx[ky+1][kx+1]
					sumY += pixel * Gy[ky+1][kx+1]
				}
			}

			// Note position of x, y
			edges[y][x] = XYToEdge(float64(sumY), float64(sumX), threshold)
		}
	}
	return edges
}
