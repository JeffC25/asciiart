package asciiart

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestToGrayScale(t *testing.T) {
	testData := []struct {
		filePath string
		width    int
	}{
		{
			filePath: filepath.Join("..", "testdata", "sample_image_0.png"),
			width:    128,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_1.png"),
			width:    128,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_2.png"),
			width:    128,
		},
		{
			filePath: filepath.Join("..", "testdata", "sample_image_3.png"),
			width:    128,
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

		pre := Preprocess(img, d.width)

		outputPath := filepath.Join("..", "testdata", "output", "TestPreproc"+strconv.Itoa(i)+".png")
		outFile, err := os.Create(outputPath)
		if err != nil {
			t.Fatalf("Failed to create TestPreproc%d.png: %v \n", i, err)
		}

		defer outFile.Close()

		err = png.Encode(outFile, pre)
		if err != nil {
			t.Fatalf("Failed to save image: %v", err)
		}

		t.Logf("Image saved as TestPreproc%d.png", i)
	}
}
