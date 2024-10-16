package asciiart

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/disintegration/gift"
)

type DoGOptions struct {
	Sigma1  float32 // Standard deviation of the first Gaussian blur
	Sigma2  float32 // Standard deviation of the second Gaussian blur
	Epsilon float32 // Threshold (0, 1) for Gaussian difference to be considered fully bright
	Tau     float32 // Emphasis of larger-scale structures
	Phi     float32 // Sharpness of edge transitions
}

func validateDoGOptions(opts DoGOptions) error {
	if opts.Epsilon > 1 || opts.Epsilon < 0 {
		return fmt.Errorf("epsilon must be between 0 and 1, inclusive")
	}

	return nil
}

// Extended thresholding function for DoG output
func tanThreshold(u, epsilon, phi float32) (float32, error) {
	if epsilon > 1 || epsilon < 0 {
		return 0, fmt.Errorf("epsilon must be between 0 and 1, inclusive")
	}
	if u >= epsilon {
		return 1, nil
	}
	return float32((1 + math.Tanh(float64(phi*(u-epsilon))))), nil
}

// Apply Difference of Gaussians as preprocessor for edge detection
func DoG(img image.Image, opts DoGOptions) (*image.Gray, error) {
	err := validateDoGOptions(opts)
	if err != nil {
		return nil, err
	}

	log.Println("Applying Difference of Gaussians preprocessing...")
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

			// Winnemoller's XDoG operator: (1 + τ) * G_1 - τ * G_2
			g1 := float32(p1.Y) / 255
			g2 := float32(p2.Y) / 255
			d, err := tanThreshold((1+opts.Tau)*g1-opts.Tau*g2, opts.Epsilon, opts.Phi)
			if err != nil {
				return nil, err
			}
			doG.Set(j, i, color.Gray{Y: uint8(math.Round(255 * float64(d)))})
		}
	}
	return doG, nil
}

func Grayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	dst := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			dst.Set(x, y, grayColor)
		}
	}
	return dst
}
