package asciiart

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestDoG(t *testing.T) {
	testData := []struct {
		filePath string
		opt      DoGOptions
	}{
		{
			filePath: filepath.Join("..", "testdata", "sample_image_0.png"),
			opt:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_1.png"),
			opt:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_2.png"),
			opt:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_3.png"),
			opt:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		},
	}

	for i, d := range testData {

		file, err := os.Open(d.filePath)
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			t.Fatalf("Failed to decode image: %v", err)
		}

		doG := DoG(img, d.opt)

		outputPath := filepath.Join("..", "testdata", "output", "TestDoG"+strconv.Itoa(i)+".png")
		outFile, err := os.Create(outputPath)
		if err != nil {
			t.Fatalf("Failed to create TestDoG%d.png: %v \n", i, err)
		}

		defer outFile.Close()

		err = png.Encode(outFile, doG)
		if err != nil {
			t.Fatalf("Failed to save image: %v", err)
		}

		t.Logf("Image saved as TestDoG%d.png", i)
	}
}

func TestXYToEdge(t *testing.T) {
	testData := []struct {
		x, y, threadhold float64
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
			res := XYToEdge(d.x, d.y, d.threadhold)
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

	edges := MapEdges(img.(*image.Gray), 0.1)

	for _, row := range edges {
		t.Log(row)
	}

}

func TestIntegration(t *testing.T) {
	testData := []struct {
		filePath   string
		dOpts      DoGOptions
		sThreshold float64
		scale      int
		eThreshold float64
	}{
		{
			filePath:   filepath.Join("..", "testdata", "sample_image_0.png"),
			dOpts:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
			sThreshold: 0.1,
			scale:      8,
			eThreshold: 0.1,
		},
		//		{
		//			filePath:   filepath.Join("..", "testdata", "sample_image_1.png"),
		//			dOpts:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		//			sThreshold: 0.1,
		//			scale:      8,
		//			eThreshold: 0.3,
		//		},
		//		{
		//			filePath:   filepath.Join("..", "testdata", "sample_image_2.png"),
		//			dOpts:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		//			sThreshold: 0.1,
		//			scale:      64,
		//			eThreshold: 0.3,
		//		},
		//		{
		//			filePath:   filepath.Join("..", "testdata", "sample_image_3.png"),
		//			dOpts:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		//			sThreshold: 0.1,
		//			scale:      64,
		//			eThreshold: 0.3,
		//		},
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

		doG := DoG(img, d.dOpts)
		edges := MapEdges(doG, d.sThreshold)
		edgesDS, _ := DownscaleEdges(edges, d.scale, d.eThreshold)

		for _, row := range edgesDS {
			t.Log(string(row))
		}

	}
}
