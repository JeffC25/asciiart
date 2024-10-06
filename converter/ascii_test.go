package asciiart

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestPreProcess(t *testing.T) {
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

func TestConvertToASCIIArt(t *testing.T) {
	testData := []struct {
		filePath string
		charset  []rune
	}{
		{
			filePath: filepath.Join("..", "testdata", "downscaled_gray_0.png"),
			charset:  []rune(" .:-=+*#%@"),
		},
		{
			filePath: filepath.Join("..", "testdata", "downscaled_gray_1.png"),
			charset:  []rune(" .:-=+*#%@"),
		},
		{
			filePath: filepath.Join("..", "testdata", "downscaled_gray_2.png"),
			charset:  []rune(" .:-=+*#%@"),
		},
		{
			filePath: filepath.Join("..", "testdata", "downscaled_gray_3.png"),
			charset:  []rune(" .:-=+*#%@"),
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

		a := ConvertToASCIIArt(img, d.charset)

		outputPath := filepath.Join("..", "testdata", "output", "ASCIIArt"+strconv.Itoa(i)+".txt")
		outFile, err := os.Create(outputPath)
		if err != nil {
			t.Fatalf("Failed to create ASCIIArt%d.txt: %v \n", i, err)
		}
		defer outFile.Close()

		os.Stdout = outFile

		PrintASCIIArt(a)

		t.Logf("Image saved as ASCIIArt%d.txt", i)
	}
}
