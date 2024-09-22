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
		sigma1   float32
		sigma2   float32
	}{
		{
			filePath: filepath.Join("..", "testdata", "sample_image_0.png"),
			sigma1:   2,
			sigma2:   4,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_1.png"),
			sigma1:   2,
			sigma2:   4,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_2.png"),
			sigma1:   2,
			sigma2:   4,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_3.png"),
			sigma1:   2,
			sigma2:   4,
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

		doG := DoG(img, d.sigma1, d.sigma2)

		outputPath := filepath.Join("..", "testdata", "output", "TestDoG"+strconv.Itoa(i)+".png")
		outFile, err := os.Create(outputPath)
		if err != nil {
			t.Fatalf("Failed to create TestDoG%d.png: %v \n", i, err)
		}

		defer outFile.Close()

		err = png.Encode(outFile, doG)
		if err != nil {
			t.Fatalf("Failed to save grayscale image: %v", err)
		}

		t.Logf("Grayscale image saved as TestDoG%d.png", i)
	}
}

func TestSobelX(t *testing.T) {
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

	edges := SobelX(img)

	outputPath := filepath.Join("..", "testdata", "output", "TestSobelX.png")
	outFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create TestSobelX.png: %v", err)
	}

	defer outFile.Close()

	err = png.Encode(outFile, edges)
	if err != nil {
		t.Fatalf("Failed to save grayscale image: %v", err)
	}

	t.Logf("Grayscale image saved as TestSobelX.png")

}

func TestSobelY(t *testing.T) {
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

	edges := SobelY(img)

	outputPath := filepath.Join("..", "testdata", "output", "TestSobelY.png")
	outFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create TestSobelY.png: %v", err)
	}

	defer outFile.Close()

	err = png.Encode(outFile, edges)
	if err != nil {
		t.Fatalf("Failed to save grayscale image: %v", err)
	}

	t.Logf("Grayscale image saved as TestSobelY.png")

}
