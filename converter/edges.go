package asciiart

import (
	"fmt"
	"image"
	"image/color"
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
	Sigma1  float32 // 1st Gaussian blur
	Sigma2  float32 // 2nd Gaussian blur
	Epsilon float32 // threshold
	Tau     float32 // Gaussian scalar coefficient
	Phi     float32 // hyperbolic tangent
}

// Extended thresholding function for DoG output
func TanhThreshold(u, epsilon, phi float32) uint8 {
	if u >= epsilon {
		return 1
	}
	return uint8((1 + math.Tanh(float64(phi*(u-epsilon)))))
}

// Apply Difference of Gaussians as preprocessor for edge detection
func DoG(img image.Image, opts DoGOptions) *image.Gray {
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
			d := TanhThreshold((1+opts.Tau)*g1-opts.Tau*g2, opts.Epsilon, opts.Phi)
			doG.Set(j, i, color.Gray{Y: 255 * d})
		}
	}
	return doG
}

// Compute angle of X Y gradients and map to discrete edges if magnitude above threshold
func XYToEdge(x, y, threshold float64) Edge {
	magnitude := math.Hypot(x, y)
	if magnitude < threshold || magnitude == 0 {
		return None
	}

	// math.Atan2 outputs -π to π radians
	angle := math.Atan2(y, x)

	// Normalize angle to the range [0, π]
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

// Map an image to a 2d slice of Edge types
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

func DownscaleEdges(edges [][]Edge, scale int, threshold float64) ([][]rune, error) {
	height := len(edges)
	width := len(edges[0])

	if height < scale || width < scale {
		return nil, fmt.Errorf("scale (%d) is larger than dimensions (%d x %d)", scale, width, height)
	}

	newHeight := height / scale
	newWidth := width / scale

	dst := make([][]rune, newHeight)
	for y := 0; y < newHeight; y++ {
		dst[y] = make([]rune, newWidth)
	}

	getSubmatrixEdge := func(x int, y int) Edge {
		edgeCounts := make(map[Edge]int)
		total := 0
		// Analyze the current submatrix of size scale x scale
		for subY := 0; subY < scale; subY++ {
			for subX := 0; subX < scale; subX++ {
				edge := edges[y*scale+subY][x*scale+subX]
				edgeCounts[edge]++
				total++
			}
		}

		if float64(edgeCounts[None])/float64(total-edgeCounts[Default]) > 1-threshold {
			return None
		}

		maxCount := 0
		maxEdge := None
		for edge, count := range edgeCounts {
			if edge != None && edge != Default && count >= maxCount {
				maxCount = count
				maxEdge = edge
			}
		}

		return maxEdge

	}

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			switch getSubmatrixEdge(x, y) {
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
