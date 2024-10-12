package asciiart

import (
	"image"
	"os"
	"path/filepath"
	"testing"
)

func TestConverter(t *testing.T) {
	testData := []string{
		filepath.Join("..", "testdata", "sample_image_0.png"),
		filepath.Join("..", "testdata", "sample_image_1.png"),
		filepath.Join("..", "testdata", "sample_image_2.png"),
		filepath.Join("..", "testdata", "sample_image_3.png"),
	}

	for _, f := range testData {
		file, err := os.Open(f)
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			t.Fatalf("Failed to decode image: %v", err)
		}

		c := NewConverter(img)
		t.Log("Converting to ascii with base and edges...")
		a, err := c.Convert()
		if err != nil {
			t.Errorf("Error converting image: %v", err)
		}

		for _, row := range a {
			t.Log(string(row))
		}

		c.DoEdges = false
		t.Log("Converting to ascii without edges...")
		a, err = c.Convert()
		if err != nil {
			t.Errorf("Error converting image: %v", err)
		}

		for _, row := range a {
			t.Log(string(row))
		}

		c.DoBase = false
		c.DoEdges = true
		t.Log("Converting to ascii with edges only...")
		a, err = c.Convert()
		if err != nil {
			t.Errorf("Error converting image: %v", err)
		}

		for _, row := range a {
			t.Log(string(row))
		}

	}
}
