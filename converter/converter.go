package asciiart

import (
	"fmt"
	"image"
)

type Converter struct {
	Img        image.Image
	CharSet    []rune
	CharWidth  int
	DOpts      DoGOptions
	SThreshold float32
	EThreshold float32
	HWeight    float32
	DoEdges    bool
	DoBase     bool
	DoDoG      bool
}

func NewConverter(img image.Image, options ...func(*Converter)) *Converter {
	c := &Converter{
		Img:        img,
		CharWidth:  175,
		CharSet:    []rune(" .:-=+*#%@"),
		DOpts:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		SThreshold: 0.15,
		EThreshold: 0.05,
		HWeight:    2.3,
		DoEdges:    true,
		DoBase:     true,
		DoDoG:      true,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func WithWidth(width int) func(*Converter) {
	return func(c *Converter) {
		c.CharWidth = width
	}
}

func WithCharset(charset []rune) func(*Converter) {
	return func(c *Converter) {
		c.CharSet = charset
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
		c.EThreshold = threshold
	}
}

func WithhWeight(hWeight float32) func(*Converter) {
	return func(c *Converter) {
		c.HWeight = hWeight
	}
}

func WithDoEdges(doEdges bool) func(*Converter) {
	return func(c *Converter) {
		c.DoEdges = doEdges
	}
}

func WithDoBase(doBase bool) func(*Converter) {
	return func(c *Converter) {
		c.DoBase = doBase
	}
}

func WithDoDoG(doDoG bool) func(*Converter) {
	return func(c *Converter) {
		c.DoDoG = doDoG
	}
}

func (c *Converter) Convert() ([][]rune, error) {
	if !c.DoEdges && !c.DoBase {
		return nil, fmt.Errorf("can not convert with both edges and base ascii disabled")
	}

	var a [][]rune
	if c.DoBase {
		fmt.Println("Mapping luminance to ascii...")
		g, err := GrayDownscale(c.Img, c.CharWidth, c.HWeight)
		if err != nil {
			return nil, err
		}

		a, err = ConvertToASCIIArt(g, c.CharSet)
		if err != nil {
			return nil, err
		}
	}

	var e [][]rune
	if c.DoEdges {
		fmt.Println("Mapping edges to ascii...")

		var d *image.Gray
		if c.DoDoG {
			var err error
			d, err = DoG(c.Img, c.DOpts)
			if err != nil {
				return nil, err
			}
		} else {
			d = image.NewGray(c.Img.Bounds())
			for y := 0; y < c.Img.Bounds().Dy(); y++ {
				for x := 0; x < c.Img.Bounds().Dx(); x++ {
					d.Set(x, y, c.Img.At(x, y))
				}
			}
		}

		m, err := MapEdges(d, c.SThreshold)
		if err != nil {
			return nil, err

		}

		e, err = DownscaleEdges(m, c.CharWidth, c.HWeight, c.EThreshold)
		if err != nil {
			return nil, err
		}
	}

	switch {
	case c.DoEdges && c.DoBase:
		dst, err := OverlayEdges(a, e)
		if err != nil {
			return nil, err
		}
		return dst, nil
	case c.DoEdges:
		return e, nil
	default:
		return a, nil
	}

}
