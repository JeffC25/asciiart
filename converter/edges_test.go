package asciiart

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestDoG(t *testing.T) {
	filePath := filepath.Join("..", "testdata", "sample_image_1.png")
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("Failed to decode image: %v", err)
	}

	doG := DoG(img, 0.5, 1.5)

	outputPath := filepath.Join("..", "testdata", "output", "TestDoG1.png")
	outFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create TestDoG1.png: %v", err)
	}

	defer outFile.Close()

	err = png.Encode(outFile, doG)
	if err != nil {
		t.Fatalf("Failed to save grayscale image: %v", err)
	}

	t.Logf("Grayscale image saved as TestDoG1.png")
}
