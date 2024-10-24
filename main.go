package main

import (
	"flag"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/jeffc25/asciiart/asciiart"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: asciiart [options] <file>")
	}

	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("Failed to open file: %v\n", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	charset := flag.String("charset", " .:-=+*#%@", "ascii characters to map image to")
	charwidth := flag.Int("width", 175, "width of ascii character conversion")
	sigma1 := flag.Float64("s1", 4, "sigma for first gaussian filter")
	sigma2 := flag.Float64("s2", 10, "sigma for second gaussian filter")
	epsilon := flag.Float64("epsilon", 0.65, "epsilon for DoG")
	tau := flag.Float64("tau", 0.8, "tau for DoG")
	phi := flag.Float64("phi", 25, "phi for DoG")
	sThreshold := flag.Float64("sthres", 0.15, "minimum threshold for sobel filter")
	eThreshold := flag.Float64("ethres", 0.05, "minimum edge density in a downscaled block")
	squash := flag.Float64("squash", 2.3, "factor to compress output height to offset aspect ratio differences between ascii and pixels")
	noEdges := flag.Bool("noedges", false, "convert without edge detection")
	noBase := flag.Bool("nobase", false, "convert without base ascii luminance mapping")
	noDoG := flag.Bool("nodog", false, "convert without Difference of Gaussians preprocessing for edge detection")

	flag.Parse()

	c := asciiart.NewConverter(img,
		asciiart.WithCharset([]rune(*charset)),
		asciiart.WithWidth(*charwidth),
		asciiart.WithDSigma1(float32(*sigma1)),
		asciiart.WithDSigma2(float32(*sigma2)),
		asciiart.WithDEpsilon(float32(*epsilon)),
		asciiart.WithDTau(float32(*tau)),
		asciiart.WithDPhi(float32(*phi)),
		asciiart.WithSThreshold(float32(*sThreshold)),
		asciiart.WithEThreshold(float32(*eThreshold)),
		asciiart.WithSquash(float32(*squash)),
		asciiart.WithDoEdges(!*noEdges),
		asciiart.WithDoBase(!*noBase),
		asciiart.WithDoDoG(!*noDoG),
	)

	a, err := c.Convert()
	if err != nil {
		log.Fatalf("Failed to convert to ascii: %v\n", err)
	}

	asciiart.PrintASCIIArt(a)
}
