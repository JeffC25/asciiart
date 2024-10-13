package asciiart

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestSquash(t *testing.T) {
	testData := []struct {
		filePath string
		scale    float32
	}{
		{
			filePath: filepath.Join("..", "testdata", "sample_image_0.png"),
			scale:    2,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_1.png"),
			scale:    2,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_2.png"),
			scale:    2,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_3.png"),
			scale:    2,
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

		pre := Squash(img, d.scale)

		outputPath := filepath.Join("..", "testdata", "output", "TestSquash"+strconv.Itoa(i)+".png")
		outFile, err := os.Create(outputPath)
		if err != nil {
			t.Fatalf("Failed to create TestSquash%d.png: %v \n", i, err)
		}

		defer outFile.Close()

		err = png.Encode(outFile, pre)
		if err != nil {
			t.Fatalf("Failed to save image: %v", err)
		}

		t.Logf("Image saved as TestSquash%d.png", i)
	}
}
