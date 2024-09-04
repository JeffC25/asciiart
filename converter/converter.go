package asciiart

import (
	"image"
)

type Converter struct {
	charset string
	scale   int
}

type Options struct {
	scale int
}

func NewConverter(opts ...func(*Converter)) *Converter {
	c := &Converter{
		scale:   8,
		charset: " .:-=+*#%@",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithScale(scale int) func(*Converter) {
	return func(c *Converter) {
		c.scale = scale
	}
}

func (c *Converter) Convert(img image.Image) string {
	// TODO

	return ""
}
