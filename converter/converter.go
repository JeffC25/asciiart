package asciiart

import (
	"fmt"
	"image"
)

type Converter struct {
	img        image.Image
	charset    []rune
	width      int
	dOpts      DoGOptions
	sThreshold float32
	eThreshold float32
	hWeight    float32
}

func NewConverter(img image.Image, options ...func(*Converter)) *Converter {
	c := &Converter{
		img:        img,
		width:      175,
		charset:    []rune(" .:-=+*#%@"),
		dOpts:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		sThreshold: 0.15,
		eThreshold: 0.05,
		hWeight:    2.3,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func WithWidth(width int) func(*Converter) {
	return func(c *Converter) {
		c.width = width
	}
}

func WithCharset(charset []rune) func(*Converter) {
	return func(c *Converter) {
		c.charset = charset
	}
}

func WithDSigma1(sigma float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Sigma1 = sigma
	}
}

func WithDSigma2(sigma float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Sigma2 = sigma
	}
}

func WithDEpsilon(epsilon float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Epsilon = epsilon
	}
}

func WithDTau(tau float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Tau = tau
	}
}

func WithDPhi(phi float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Phi = phi
	}
}

func WithSThreshold(threshold float32) func(*Converter) {
	return func(c *Converter) {
		c.sThreshold = threshold
	}
}

func WithEThreshold(threshold float32) func(*Converter) {
	return func(c *Converter) {
		c.sThreshold = threshold
	}
}

func WithhWeight(hWeight float32) func(*Converter) {
	return func(c *Converter) {
		c.hWeight = hWeight
	}
}

func (c *Converter) Convert() ([][]rune, error) {
	fmt.Println("Mapping luminance to ascii...")
	a := ConvertToASCIIArt(GrayDownscale(c.img, c.width, c.hWeight), c.charset)

	fmt.Println("Mapping edges to ascii...")
	e, err := DownscaleEdges(MapEdges(DoG(c.img, c.dOpts), c.sThreshold), c.width, c.hWeight, c.eThreshold)
	if err != nil {
		return nil, err
	}

	dst, err := OverlayEdges(a, e)
	if err != nil {
		return nil, err
	}

	return dst, nil
}
