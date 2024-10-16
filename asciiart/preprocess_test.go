package asciiart

import (
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
			opt:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.85, Phi: 15},
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_1.png"),
			opt:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.85, Phi: 15},
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_2.png"),
			opt:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.85, Phi: 15},
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_3.png"),
			opt:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.85, Phi: 15},
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

		doG, err := DoG(img, d.opt)
		if err != nil {
			t.Fatalf("Failed to perform Difference of Gaussians: %v\n", err)
		}

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
