package asciiart

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestXYToEdge(t *testing.T) {
	testData := []struct {
		x, y, threadhold float32
		want             Edge
	}{
		{0, 0, 0, None},
		{2, 2, 8, None},
		{2, 0, 1, Horizontal},
		{-2, 0, 1, Horizontal},
		{2, 2, 1, DiagonalUp},
		{-2, -2, 1, DiagonalUp},
		{0, 2, 1, Vertical},
		{0, -2, 1, Vertical},
		{2, -2, 1, DiagonalDown},
		{-2, 2, 1, DiagonalDown},
	}

	for _, d := range testData {
		testname := fmt.Sprintf("%.2f,%.2f,%.2f", d.x, d.y, d.threadhold)
		t.Run(testname, func(t *testing.T) {
			res := xyToEdge(d.x, d.y, d.threadhold)
			if res != d.want {
				t.Errorf("got %d, want %d", res, d.want)
			}
		})
	}
}

func TestMapEdges(t *testing.T) {
	filePath := filepath.Join("..", "testdata", "downscaled_gray_0.png")
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("Failed to decode image: %v", err)
	}

	edges, err := MapEdges(img.(*image.Gray), 0.1)
	if err != nil {
		t.Fatalf("Failed to map edges: %v\n", err)
	}

	for _, row := range edges {
		t.Log(row)
	}

}

func TestEdgesIntegration(t *testing.T) {
	testData := []struct {
		filePath   string
		dOpts      DoGOptions
		newWidth   int
		sThreshold float32
		eThreshold float32
	}{
		{
			filePath:   filepath.Join("..", "testdata", "sample_image_0.png"),
			dOpts:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.8, Phi: 25},
			newWidth:   185,
			sThreshold: 0.15,
			eThreshold: 0.05,
		},
		{
			filePath:   filepath.Join("..", "testdata", "sample_image_1.png"),
			dOpts:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.8, Phi: 25},
			newWidth:   185,
			sThreshold: 0.15,
			eThreshold: 0.05,
		},
		{
			filePath:   filepath.Join("..", "testdata", "sample_image_2.png"),
			dOpts:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.85, Phi: 25},
			newWidth:   185,
			sThreshold: 0.15,
			eThreshold: 0.05,
		},
		{
			filePath:   filepath.Join("..", "testdata", "sample_image_3.png"),
			dOpts:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.85, Phi: 25},
			newWidth:   185,
			sThreshold: 0.15,
			eThreshold: 0.05,
		},
	}

	for _, d := range testData {

		file, err := os.Open(d.filePath)
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			t.Fatalf("Failed to decode image: %v", err)
		}

		doG, err := DoG(img, d.dOpts)
		if err != nil {
			t.Fatalf("Failed to perform Difference of Gaussians: %v\n", err)
		}

		edges, err := MapEdges(doG, d.sThreshold)
		if err != nil {
			t.Fatalf("Failed to map edges: %v\n", err)
		}

		edgesDS, _ := DownscaleEdges(edges, d.newWidth, 2.3, d.eThreshold)
		for _, row := range edgesDS {
			t.Log(string(row))
		}

	}
}
