package asciiart

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/disintegration/gift"
)

type Edge int

const (
	Default      Edge = iota
	None              // Non-edge - differentiate from default/undetermined
	Horizontal        // "_"
	Vertical          // "|"
	DiagonalUp        // "/"
	DiagonalDown      // "\"
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

	log.Println("Applying Difference of Gaussians")
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

// Compute angle of X Y gradients and map to discrete edges if magnitude above threshold
func xyToEdge(x, y, threshold float32) Edge {
	magnitude := math.Hypot(float64(y), float64(x))
	if magnitude < float64(threshold) || magnitude == 0 {
		return None
	}

	// math.Atan2 outputs -π to π radians
	angle := math.Atan2(float64(y), float64(x))

	// Normalize angle to the range [0, π]
	angle = math.Mod(angle, 2*math.Pi)
	if angle < 0 {
		angle += math.Pi
	}

	// Map the angle to the appropriate edge type
	switch {
	case angle >= 0 && angle < math.Pi/8 || angle >= 15*math.Pi/8 && angle <= 2*math.Pi:
		return Horizontal
	case angle >= math.Pi/8 && angle < 3*math.Pi/8:
		return DiagonalUp
	case angle >= 3*math.Pi/8 && angle < 5*math.Pi/8:
		return Vertical
	case angle >= 5*math.Pi/8 && angle < 7*math.Pi/8:
		return DiagonalDown
	default:
		return Horizontal // Handles angles close to multiples of π
	}
}

// Map an image to a 2d slice of Edge types
func MapEdges(img *image.Gray, sobelThreshold float32) ([][]Edge, error) {
	if sobelThreshold < 0 || sobelThreshold > 1 {
		return nil, fmt.Errorf("sobel filter threshold must be between 0 and 1, inclusive")
	}
	log.Println("Mapping edges...")
	threshold := sobelThreshold * float32(math.Hypot(255*4, 255*4))

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
			edges[y][x] = xyToEdge(float32(sumY), float32(sumX), threshold)
		}
	}
	return edges, nil
}

func DownscaleEdges(edges [][]Edge, newWidth int, hWeight, threshold float32) ([][]rune, error) {
	if newWidth <= 0 {
		return nil, fmt.Errorf("non-positive newWidth: %d", newWidth)
	}

	if hWeight <= 0 {
		return nil, fmt.Errorf("non-positive hWeight: %2f", hWeight)
	}

	if threshold < 0 || threshold > 1 {
		return nil, fmt.Errorf("threshold needs to be between 0 and 1: %2f", threshold)
	}

	log.Println("Downscaling edges...")
	height := len(edges)
	width := len(edges[0])

	xScale := float64(width) / float64(newWidth)
	yScale := xScale * float64(hWeight)
	newHeight := int(math.Floor(float64(height) / yScale))

	dst := make([][]rune, newHeight)
	for y := 0; y < newHeight; y++ {
		dst[y] = make([]rune, newWidth)
	}

	getSubmatrixEdge := func(x int, y int) (Edge, error) {
		edgeCounts := make(map[Edge]int)
		total := 0
		// Analyze the current submatrix of size scale x scale
		for subY := 0; float64(subY) < yScale; subY++ {
			for subX := 0; float64(subX) < xScale; subX++ {
				i := int(math.Floor(float64(y)*yScale)) + subY
				j := int(math.Floor(float64(x)*xScale)) + subX

				if i >= len(edges) {
					return None, fmt.Errorf("y out of range: %d from %d", i, len(edges))
				}
				if j >= len(edges[0]) {
					return None, fmt.Errorf("x out of range: %d from %d", j, len(edges[0]))
				}

				edge := edges[i][j]
				edgeCounts[edge]++
				total++
			}
		}

		maxCount := 0
		maxEdge := None
		for edge, count := range edgeCounts {
			if edge != None && edge != Default && count >= maxCount {
				maxCount = count
				maxEdge = edge
			}
		}

		if float32(maxCount)/float32(total-edgeCounts[Default]) > threshold {
			return maxEdge, nil
		}
		return None, nil
	}

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			e, err := getSubmatrixEdge(x, y)
			if err != nil {
				return nil, err
			}
			switch e {
			case Horizontal:
				dst[y][x] = '_'
			case DiagonalUp:
				dst[y][x] = '/'
			case Vertical:
				dst[y][x] = '|'
			case DiagonalDown:
				dst[y][x] = '\\'
			default:
				dst[y][x] = ' '
			}
		}
	}

	return dst, nil
}

func OverlayEdges(base, edges [][]rune) ([][]rune, error) {
	log.Println("Overlaying edges...")
	width := len(base[0])
	height := len(base)

	if width != len(edges[0]) || height != len(edges) {
		return nil, fmt.Errorf("mismatched dimensions: %d x %d and %d x %d", width, height, len(edges), len(edges[0]))
	}

	dst := make([][]rune, height)
	for y := 0; y < height; y++ {
		dst[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			if edges[y][x] == ' ' {
				dst[y][x] = base[y][x]
			} else {
				dst[y][x] = edges[y][x]
			}
		}
	}

	return dst, nil
}
