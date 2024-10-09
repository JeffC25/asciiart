package asciiart

import (
	"fmt"
	"image"
)

type ynoOpt int8

const (
	Yes ynoOpt = iota
	No
	Only
)

type Converter struct {
	Img        image.Image
	Charset    []rune
	Width      int
	DOpts      DoGOptions
	SThreshold float32
	EThreshold float32
	HWeight    float32
	IncEdges   ynoOpt
}

func NewConverter(img image.Image, options ...func(*Converter)) *Converter {
	c := &Converter{
		Img:        img,
		Width:      175,
		Charset:    []rune(" .:-=+*#%@"),
		DOpts:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		SThreshold: 0.15,
		EThreshold: 0.05,
		HWeight:    2.3,
		IncEdges:   Yes,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func WithWidth(width int) func(*Converter) {
	return func(c *Converter) {
		c.Width = width
	}
}

func WithCharset(charset []rune) func(*Converter) {
	return func(c *Converter) {
		c.Charset = charset
	}
}

func WithDSigma1(sigma float32) func(*Converter) {
	return func(c *Converter) {
		c.DOpts.Sigma1 = sigma
	}
}

func WithDSigma2(sigma float32) func(*Converter) {
	return func(c *Converter) {
		c.DOpts.Sigma2 = sigma
	}
}

func WithDEpsilon(epsilon float32) func(*Converter) {
	return func(c *Converter) {
		c.DOpts.Epsilon = epsilon
	}
}

func WithDTau(tau float32) func(*Converter) {
	return func(c *Converter) {
		c.DOpts.Tau = tau
	}
}

func WithDPhi(phi float32) func(*Converter) {
	return func(c *Converter) {
		c.DOpts.Phi = phi
	}
}

func WithSThreshold(threshold float32) func(*Converter) {
	return func(c *Converter) {
		c.SThreshold = threshold
	}
}

func WithEThreshold(threshold float32) func(*Converter) {
	return func(c *Converter) {
		c.SThreshold = threshold
	}
}

func WithhWeight(hWeight float32) func(*Converter) {
	return func(c *Converter) {
		c.HWeight = hWeight
	}
}

func WithEdges(yno ynoOpt) func(*Converter) {
	return func(c *Converter) {
		c.IncEdges = yno
	}
}

func (c *Converter) Convert() ([][]rune, error) {
	var a [][]rune
	if c.IncEdges != Only {
		fmt.Println("Mapping luminance to ascii...")
		g, err := (GrayDownscale(c.Img, c.Width, c.HWeight))
		if err != nil {
			return nil, err
		}

		a, err = ConvertToASCIIArt(g, c.Charset)
		if err != nil {
			return nil, err
		}
	}

	var e [][]rune
	if c.IncEdges != No {
		fmt.Println("Mapping edges to ascii...")
		d, err := DoG(c.Img, c.DOpts)
		if err != nil {
			return nil, err
		}

		m, err := MapEdges(d, c.SThreshold)
		if err != nil {
			return nil, err

		}

		e, err = DownscaleEdges(m, c.Width, c.HWeight, c.EThreshold)
		if err != nil {
			return nil, err
		}
	}

	switch c.IncEdges {
	case Only:
		return e, nil
	case Yes:
		dst, err := OverlayEdges(a, e)
		if err != nil {
			return nil, err

		}
		return dst, nil
	default: // no edges
		return a, nil
	}

}
