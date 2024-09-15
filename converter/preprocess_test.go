package asciiart

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestToGrayScale(t *testing.T) {
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

	grayImg := ToGrayscale(img)

	outputPath := filepath.Join("..", "testdata", "output", "TestToGrayScale1.png")
	outFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create TestToGrayScale1.png: %v", err)
	}

	defer outFile.Close()

	err = png.Encode(outFile, grayImg)
	if err != nil {
		t.Fatalf("Failed to save grayscale image: %v", err)
	}

	t.Logf("Grayscale image saved as TestToGrayScale1.png")
}
