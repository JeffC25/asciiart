package asciiart

import (
	"fmt"
	"image"
	"log"
)

type Converter struct {
	Img        image.Image // image to convert
	CharSet    []rune      // ascii characters to map image to
	CharWidth  int         // width of ascii character conversion
	DOpts      DoGOptions  // options for Difference of Gaussians preprocessing for edge detection
	SThreshold float32     // minimum threshold (0 to 1) for sobel filter to recognize gradient as edge
	EThreshold float32     // minimum edge density (0 to 1) in a downscaled block for it to be considered an edge
	Squash     float32     // factor to compress output height to offset aspect ratio differences between ascii and pixels
	DoEdges    bool        // whether to apply edge detection
	DoBase     bool        // whether to apply base ascii luminance mapping
	DoDoG      bool        // whether to apply Difference of Gaussians preprocessing for edge detection
}

func NewConverter(img image.Image, options ...func(*Converter)) *Converter {
	c := &Converter{
		Img:        img,
		CharWidth:  175,
		CharSet:    []rune(" .:-=+*#%@"),
		DOpts:      DoGOptions{Sigma1: 4, Sigma2: 10, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		SThreshold: 0.15,
		EThreshold: 0.05,
		Squash:     2.3,
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

func WithSquash(squash float32) func(*Converter) {
	return func(c *Converter) {
		c.Squash = squash
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
		return nil, fmt.Errorf("both edge detection and base ASCII generation are disabled; please enable at least one option")
	}

	var a [][]rune
	if c.DoBase {
		log.Println("Mapping luminance to ascii...")
		g, err := GrayDownscale(c.Img, c.CharWidth, c.Squash)
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
		log.Println("Mapping edges to ascii...")

		var d *image.Gray
		if c.DoDoG {
			var err error
			d, err = DoG(c.Img, c.DOpts)
			if err != nil {
				return nil, err
			}
		} else {
			d = Grayscale(c.Img)
		}

		m, err := MapEdges(d, c.SThreshold)
		if err != nil {
			return nil, err

		}

		e, err = DownscaleEdges(m, c.CharWidth, c.Squash, c.EThreshold)
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
